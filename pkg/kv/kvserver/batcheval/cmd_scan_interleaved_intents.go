// Copyright 2021 The Cockroach Authors.
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
	"time"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/spanset"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/storage/enginepb"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/cockroachdb/errors"
)

func init() {
	RegisterReadOnlyCommand(roachpb.ScanInterleavedIntents, declareKeysScanInterleavedIntents, ScanInterleavedIntents)
}

func declareKeysScanInterleavedIntents(
	rs ImmutableRangeState,
	_ *roachpb.Header,
	_ roachpb.Request,
	latchSpans, _ *spanset.SpanSet,
	_ time.Duration,
) {
	latchSpans.AddNonMVCC(spanset.SpanReadOnly, roachpb.Span{Key: keys.RangeDescriptorKey(rs.GetStartKey())})
}

// ScanInterleavedIntents returns intents encountered in the provided span.
// These intents are then resolved in the separated intents migration, the
// usual caller for this request.
func ScanInterleavedIntents(
	ctx context.Context, reader storage.Reader, cArgs CommandArgs, response roachpb.Response,
) (result.Result, error) {
	req := cArgs.Args.(*roachpb.ScanInterleavedIntentsRequest)
	resp := response.(*roachpb.ScanInterleavedIntentsResponse)

	// Put a limit on memory usage by scanning for at least maxIntentCount
	// intents or maxIntentBytes in intent values, whichever is reached first,
	// then returning those.
	const maxIntentCount = 1000
	const maxIntentBytes = 1 << 20 // 1MB
	iter := reader.NewEngineIterator(storage.IterOptions{
		LowerBound: req.Key,
		UpperBound: req.EndKey,
	})
	defer iter.Close()
	valid, err := iter.SeekEngineKeyGE(storage.EngineKey{Key: req.Key})
	intentCount := 0
	intentBytes := 0

	for ; valid && err == nil; valid, err = iter.NextEngineKey() {
		key, err := iter.EngineKey()
		if err != nil {
			return result.Result{}, err
		}
		if !key.IsMVCCKey() {
			// This should never happen, as the only non-MVCC keys are lock table
			// keys and those are in the local keyspace. Return an error.
			return result.Result{}, errors.New("encountered non-MVCC key during lock table migration")
		}
		mvccKey, err := key.ToMVCCKey()
		if err != nil {
			return result.Result{}, err
		}
		if !mvccKey.Timestamp.IsEmpty() {
			// Versioned value - not an intent.
			//
			// TODO(bilal): Explore seeking here in case there are keys with lots of
			// versioned values.
			continue
		}

		val := iter.Value()
		meta := enginepb.MVCCMetadata{}
		if err := protoutil.Unmarshal(val, &meta); err != nil {
			return result.Result{}, err
		}
		if meta.IsInline() {
			// Inlined value - not an intent.
			continue
		}

		if intentCount >= maxIntentCount || intentBytes >= maxIntentBytes {
			// Batch limit reached - cut short this batch here. This kv
			// will be added to txnIntents on the next iteration of the outer loop.
			resp.ResumeSpan = &roachpb.Span{
				Key:    mvccKey.Key,
				EndKey: req.EndKey,
			}
			break
		}
		resp.Intents = append(resp.Intents, roachpb.MakeIntent(meta.Txn, mvccKey.Key))
		intentCount++
		intentBytes += len(val)
	}

	return result.Result{}, nil
}
