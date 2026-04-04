// Copyright 2016 The Cockroach Authors.
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

package kvserver_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"
)

func TestComputeStatsForKeySpan(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	serv, _, _ := serverutils.StartServer(t, base.TestServerArgs{})
	s := serv.(*server.TestServer)
	defer s.Stopper().Stop(ctx)
	store, err := s.Stores().GetStore(s.GetFirstStoreID())
	require.NoError(t, err)
	// Create a number of ranges using splits.
	splitKeys := []string{"a", "c", "e", "g", "i"}
	for _, k := range splitKeys {
		_, _, err := s.SplitRange(roachpb.Key(k))
		require.NoError(t, err)
	}

	// Wait for splits to finish.
	testutils.SucceedsSoon(t, func() error {
		repl := store.LookupReplica(roachpb.RKey("z"))
		if actualRSpan := repl.Desc().RSpan(); !actualRSpan.Key.Equal(roachpb.RKey("i")) {
			return errors.Errorf("expected range %s to begin at key 'i'", repl)
		}
		return nil
	})

	// Create some keys across the ranges.
	incKeys := []string{"b", "bb", "bbb", "d", "dd", "h"}
	for _, k := range incKeys {
		if _, err := store.DB().Inc(context.Background(), []byte(k), 5); err != nil {
			t.Fatal(err)
		}
	}

	// Verify stats across different spans.
	for _, tcase := range []struct {
		startKey       string
		endKey         string
		expectedRanges int
		expectedKeys   int64
	}{
		{"a", "i", 4, 6},
		{"a", "c", 1, 3},
		{"b", "e", 2, 5},
		{"e", "i", 2, 1},
	} {
		start, end := tcase.startKey, tcase.endKey
		result, err := store.ComputeStatsForKeySpan(
			roachpb.RKey(start), roachpb.RKey(end))
		if err != nil {
			t.Fatal(err)
		}
		if a, e := result.ReplicaCount, tcase.expectedRanges; a != e {
			t.Errorf("Expected %d ranges in span [%s - %s], found %d", e, start, end, a)
		}
		if a, e := result.MVCC.LiveCount, tcase.expectedKeys; a != e {
			t.Errorf("Expected %d keys in span [%s - %s], found %d", e, start, end, a)
		}
	}
}
