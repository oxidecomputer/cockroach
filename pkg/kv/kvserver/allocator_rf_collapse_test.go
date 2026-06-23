// Copyright 2014 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package kvserver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness/livenesspb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// Test demonstrating the behaviors that led to omicron#10658.
func TestAllocatorDownReplicatesOnColdLivenessCache(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()

	var coldLiveness *liveness.NodeLiveness
	stopper, g, mc, sp, _ := createTestStorePool(ctx,
		TestTimeUntilStoreDeadOff, false, /* deterministic */
		func() int {
			if coldLiveness == nil {
				return 0 // not reached in this test; nodeCountFn is unused before construction
			}
			return coldLiveness.GetNodeCount()
		},
		livenesspb.NodeLivenessStatus_LIVE)
	defer stopper.Stop(ctx)

	// We don't Start() the NodeLiveness, since we are about only the cache and
	// GetNodeCount(), not the heartbeat loop.
	clock := hlc.NewClock(mc.UnixNano, time.Nanosecond)
	coldLiveness = liveness.NewNodeLiveness(liveness.NodeLivenessOptions{
		AmbientCtx:              log.MakeTestingAmbientContext(stopper.Tracer()),
		Stopper:                 stopper,
		Settings:                cluster.MakeTestingClusterSettings(),
		Gossip:                  g,
		Clock:                   clock,
		LivenessThreshold:       time.Minute,
		RenewalDuration:         time.Second,
		HistogramWindowInterval: time.Minute,
	})

	a := MakeAllocator(sp, func(string) (time.Duration, bool) {
		return 0, true
	}, nil /* knobs */, nil /* storeMetrics */)

	// RF 5, with 5 healthy replicas on 5 live stores.
	conf := roachpb.SpanConfig{NumReplicas: 5}
	allFive := []roachpb.StoreID{1, 2, 3, 4, 5}

	testCases := []struct {
		cacheRecords        int
		expectedNumReplicas int
		expectedAction      AllocatorAction
	}{
		// < 5 leads to an effective RF of 3. (These are in ascending order to
		// make testing easier.)
		{cacheRecords: 0, expectedNumReplicas: 3, expectedAction: AllocatorRemoveVoter},
		{cacheRecords: 1, expectedNumReplicas: 3, expectedAction: AllocatorRemoveVoter},
		{cacheRecords: 2, expectedNumReplicas: 3, expectedAction: AllocatorRemoveVoter},
		{cacheRecords: 3, expectedNumReplicas: 3, expectedAction: AllocatorRemoveVoter},
		{cacheRecords: 4, expectedNumReplicas: 3, expectedAction: AllocatorRemoveVoter},
		// 5 leads to an effective RF of 5.
		{cacheRecords: 5, expectedNumReplicas: 5, expectedAction: AllocatorConsiderRebalance},
	}

	nextNode := roachpb.NodeID(1)
	for _, c := range testCases {
		t.Run(fmt.Sprintf("cacheRecords=%d", c.cacheRecords), func(t *testing.T) {
			// Grow the real cache through maybeUpdate.
			for coldLiveness.GetNodeCount() < c.cacheRecords {
				coldLiveness.TestingMaybeUpdate(ctx, liveness.Record{
					Liveness: livenesspb.Liveness{
						NodeID:     nextNode,
						Epoch:      1,
						Membership: livenesspb.MembershipStatus_ACTIVE,
					},
				})
				nextNode++
			}
			require.Equalf(t, c.cacheRecords, coldLiveness.GetNodeCount(),
				"the real liveness cache should hold exactly %d of 5 records", c.cacheRecords)

			mockStorePool(sp, allFive, nil, nil, nil, nil, nil)
			desc := makeDescriptor(allFive)

			clusterNodes := a.storePool.ClusterNodeCount()
			require.Equalf(t, c.cacheRecords, clusterNodes,
				"ClusterNodeCount() should flow from the real NodeLiveness.GetNodeCount() (%d records)",
				c.cacheRecords)

			effectiveNumReplicas := GetNeededVoters(conf.NumReplicas, clusterNodes)
			require.Equalf(t, c.expectedNumReplicas, effectiveNumReplicas,
				"GetNeededVoters(5, %d) sizes a healthy RF-5 range to %d replicas",
				clusterNodes, effectiveNumReplicas)

			action, _ := a.ComputeAction(ctx, conf, &desc)
			require.Equalf(t, c.expectedAction.String(), action.String(),
				"5 live healthy replicas + a liveness cache of %d/5 records -> allocator action %s",
				c.cacheRecords, action)
		})
	}
}
