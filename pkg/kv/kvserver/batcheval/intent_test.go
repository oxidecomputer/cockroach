// Copyright 2020 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package batcheval

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// instrumentedEngine wraps a storage.Engine and allows for various methods in
// the interface to be instrumented for testing purposes.
type instrumentedEngine struct {
	storage.Engine

	onNewIterator func(storage.IterOptions)
	// ... can be extended ...
}

func (ie *instrumentedEngine) NewMVCCIterator(
	iterKind storage.MVCCIterKind, opts storage.IterOptions,
) storage.MVCCIterator {
	if ie.onNewIterator != nil {
		ie.onNewIterator(opts)
	}
	return ie.Engine.NewMVCCIterator(iterKind, opts)
}

// TestCollectIntentsUsesSameIterator tests that all uses of CollectIntents
// (currently only by READ_UNCOMMITTED Gets, Scans, and ReverseScans) use the
// same cached iterator (prefix or non-prefix) for their initial read and their
// provisional value collection for any intents they find.
func TestCollectIntentsUsesSameIterator(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	key := roachpb.Key("key")
	ts := hlc.Timestamp{WallTime: 123}
	header := roachpb.Header{
		Timestamp:       ts,
		ReadConsistency: roachpb.READ_UNCOMMITTED,
	}
	evalCtx := (&MockEvalCtx{ClusterSettings: cluster.MakeTestingClusterSettings()}).EvalContext()

	testCases := []struct {
		name              string
		run               func(*testing.T, storage.ReadWriter) (intents []roachpb.KeyValue, _ error)
		expPrefixIters    int
		expNonPrefixIters int
	}{
		{
			name: "get",
			run: func(t *testing.T, db storage.ReadWriter) ([]roachpb.KeyValue, error) {
				req := &roachpb.GetRequest{
					RequestHeader: roachpb.RequestHeader{Key: key},
				}
				var resp roachpb.GetResponse
				if _, err := Get(ctx, db, CommandArgs{Args: req, Header: header, EvalCtx: evalCtx}, &resp); err != nil {
					return nil, err
				}
				if resp.IntentValue == nil {
					return nil, nil
				}
				return []roachpb.KeyValue{{Key: key, Value: *resp.IntentValue}}, nil
			},
			expPrefixIters:    2,
			expNonPrefixIters: 0,
		},
		{
			name: "scan",
			run: func(t *testing.T, db storage.ReadWriter) ([]roachpb.KeyValue, error) {
				req := &roachpb.ScanRequest{
					RequestHeader: roachpb.RequestHeader{Key: key, EndKey: key.Next()},
				}
				var resp roachpb.ScanResponse
				if _, err := Scan(ctx, db, CommandArgs{Args: req, Header: header, EvalCtx: evalCtx}, &resp); err != nil {
					return nil, err
				}
				return resp.IntentRows, nil
			},
			expPrefixIters:    0,
			expNonPrefixIters: 2,
		},
		{
			name: "reverse scan",
			run: func(t *testing.T, db storage.ReadWriter) ([]roachpb.KeyValue, error) {
				req := &roachpb.ReverseScanRequest{
					RequestHeader: roachpb.RequestHeader{Key: key, EndKey: key.Next()},
				}
				var resp roachpb.ReverseScanResponse
				if _, err := ReverseScan(ctx, db, CommandArgs{Args: req, Header: header, EvalCtx: evalCtx}, &resp); err != nil {
					return nil, err
				}
				return resp.IntentRows, nil
			},
			expPrefixIters:    0,
			expNonPrefixIters: 2,
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			// Test with and without deletion intents. If a READ_UNCOMMITTED request
			// encounters an intent whose provisional value is a deletion tombstone,
			// the request should ignore the intent and should not return any
			// corresponding intent row.
			testutils.RunTrueAndFalse(t, "deletion intent", func(t *testing.T, delete bool) {
				db := &instrumentedEngine{Engine: storage.NewDefaultInMemForTesting()}
				defer db.Close()

				// Write an intent.
				val := roachpb.MakeValueFromBytes([]byte("val"))
				txn := roachpb.MakeTransaction("test", key, roachpb.NormalUserPriority, ts, 0, 1)
				var err error
				if delete {
					err = storage.MVCCDelete(ctx, db, nil, key, ts, &txn)
				} else {
					err = storage.MVCCPut(ctx, db, nil, key, ts, val, &txn)
				}
				require.NoError(t, err)

				// Instrument iterator creation, count prefix vs. non-prefix iters.
				var prefixIters, nonPrefixIters int
				db.onNewIterator = func(opts storage.IterOptions) {
					if opts.Prefix {
						prefixIters++
					} else {
						nonPrefixIters++
					}
				}

				intents, err := c.run(t, db)
				require.NoError(t, err)

				// Assert proper intent values.
				if delete {
					require.Len(t, intents, 0)
				} else {
					expIntentVal := val
					expIntentVal.Timestamp = ts
					expIntentKeyVal := roachpb.KeyValue{Key: key, Value: expIntentVal}
					require.Len(t, intents, 1)
					require.Equal(t, expIntentKeyVal, intents[0])
				}

				// Assert proper iterator use.
				require.Equal(t, c.expPrefixIters, prefixIters)
				require.Equal(t, c.expNonPrefixIters, nonPrefixIters)
				require.Equal(t, c.expNonPrefixIters, nonPrefixIters)
			})
		})
	}
}
