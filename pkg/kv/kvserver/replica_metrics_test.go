// Copyright 2020 The Cockroach Authors.
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
	"testing"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvserverpb"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness/livenesspb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestCalcRangeCounterIsLiveMap(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	leaseStatus := kvserverpb.LeaseStatus{
		Lease: roachpb.Lease{
			Replica: roachpb.ReplicaDescriptor{
				NodeID:  10,
				StoreID: 11,
			},
		},
		State: kvserverpb.LeaseState_VALID,
	}

	// Regression test for a bug, see:
	// https://github.com/cockroachdb/cockroach/pull/39936#pullrequestreview-359059629

	threeVotersAndSingleNonVoter := roachpb.NewRangeDescriptor(123, roachpb.RKeyMin, roachpb.RKeyMax,
		roachpb.MakeReplicaSet([]roachpb.ReplicaDescriptor{
			{NodeID: 10, StoreID: 11, ReplicaID: 12, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 100, StoreID: 110, ReplicaID: 120, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 1000, StoreID: 1100, ReplicaID: 1200, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 2000, StoreID: 2100, ReplicaID: 2200, Type: roachpb.ReplicaTypeNonVoter()},
		}))

	oneVoterAndThreeNonVoters := roachpb.NewRangeDescriptor(123, roachpb.RKeyMin, roachpb.RKeyMax,
		roachpb.MakeReplicaSet([]roachpb.ReplicaDescriptor{
			{NodeID: 10, StoreID: 11, ReplicaID: 12, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 100, StoreID: 110, ReplicaID: 120, Type: roachpb.ReplicaTypeNonVoter()},
			{NodeID: 1000, StoreID: 1100, ReplicaID: 1200, Type: roachpb.ReplicaTypeNonVoter()},
			{NodeID: 2000, StoreID: 2100, ReplicaID: 2200, Type: roachpb.ReplicaTypeNonVoter()},
		}))

	{
		ctr, down, under, over := calcRangeCounter(1100, threeVotersAndSingleNonVoter, leaseStatus, liveness.IsLiveMap{
			1000: liveness.IsLiveMapEntry{IsLive: true}, // by NodeID
		}, 3 /* numVoters */, 4 /* numReplicas */, 4 /* clusterNodes */)

		require.True(t, ctr)
		require.True(t, down)
		require.True(t, under)
		require.False(t, over)
	}

	{
		ctr, down, under, over := calcRangeCounter(1000, threeVotersAndSingleNonVoter, leaseStatus, liveness.IsLiveMap{
			1000: liveness.IsLiveMapEntry{IsLive: false},
		}, 3 /* numVoters */, 4 /* numReplicas */, 4 /* clusterNodes */)

		// Does not confuse a non-live entry for a live one. In other words,
		// does not think that the liveness map has only entries for live nodes.
		require.False(t, ctr)
		require.False(t, down)
		require.False(t, under)
		require.False(t, over)
	}

	{
		ctr, down, under, over := calcRangeCounter(11, threeVotersAndSingleNonVoter, leaseStatus, liveness.IsLiveMap{
			10:   liveness.IsLiveMapEntry{IsLive: true},
			100:  liveness.IsLiveMapEntry{IsLive: true},
			1000: liveness.IsLiveMapEntry{IsLive: true},
			2000: liveness.IsLiveMapEntry{IsLive: true},
		}, 3 /* numVoters */, 4 /* numReplicas */, 4 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.False(t, under)
		require.False(t, over)
	}

	{
		// Single non-voter dead
		ctr, down, under, over := calcRangeCounter(11, oneVoterAndThreeNonVoters, leaseStatus, liveness.IsLiveMap{
			10:   liveness.IsLiveMapEntry{IsLive: true},
			100:  liveness.IsLiveMapEntry{IsLive: true},
			1000: liveness.IsLiveMapEntry{IsLive: false},
			2000: liveness.IsLiveMapEntry{IsLive: true},
		}, 1 /* numVoters */, 4 /* numReplicas */, 4 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.True(t, under)
		require.False(t, over)
	}

	{
		// All non-voters are dead, but range is not unavailable
		ctr, down, under, over := calcRangeCounter(11, oneVoterAndThreeNonVoters, leaseStatus, liveness.IsLiveMap{
			10:   liveness.IsLiveMapEntry{IsLive: true},
			100:  liveness.IsLiveMapEntry{IsLive: false},
			1000: liveness.IsLiveMapEntry{IsLive: false},
			2000: liveness.IsLiveMapEntry{IsLive: false},
		}, 1 /* numVoters */, 4 /* numReplicas */, 4 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.True(t, under)
		require.False(t, over)
	}

	{
		// More non-voters than needed
		ctr, down, under, over := calcRangeCounter(11, oneVoterAndThreeNonVoters, leaseStatus, liveness.IsLiveMap{
			10:   liveness.IsLiveMapEntry{IsLive: true},
			100:  liveness.IsLiveMapEntry{IsLive: true},
			1000: liveness.IsLiveMapEntry{IsLive: true},
			2000: liveness.IsLiveMapEntry{IsLive: true},
		}, 1 /* numVoters */, 3 /* numReplicas */, 4 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.False(t, under)
		require.True(t, over)
	}
}

// TODO-RAINCLAUDE: omicron#10658 — the policy floor applied to the gauges, so
// they agree with the allocator. Covers: the madrid phantom window (healthy
// range no longer reads over-replicated; a range whose voters the leaseholder
// cannot account for reads under-replicated instead of quiet), the
// dead-decommission and even-steady states (no longer permanently
// over-replicated), and genuine over-replication versus the configured RF
// (still detected — the floor is capped at numVoters).
func TestCalcRangeCounterPolicyFloor(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	leaseStatus := kvserverpb.LeaseStatus{
		Lease: roachpb.Lease{
			Replica: roachpb.ReplicaDescriptor{
				NodeID:  1,
				StoreID: 10,
			},
		},
		State: kvserverpb.LeaseState_VALID,
	}

	fiveVoters := roachpb.NewRangeDescriptor(123, roachpb.RKeyMin, roachpb.RKeyMax,
		roachpb.MakeReplicaSet([]roachpb.ReplicaDescriptor{
			{NodeID: 1, StoreID: 10, ReplicaID: 1, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 2, StoreID: 20, ReplicaID: 2, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 3, StoreID: 30, ReplicaID: 3, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 4, StoreID: 40, ReplicaID: 4, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 5, StoreID: 50, ReplicaID: 5, Type: roachpb.ReplicaTypeVoterFull()},
		}))

	fourVoters := roachpb.NewRangeDescriptor(124, roachpb.RKeyMin, roachpb.RKeyMax,
		roachpb.MakeReplicaSet([]roachpb.ReplicaDescriptor{
			{NodeID: 1, StoreID: 10, ReplicaID: 1, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 2, StoreID: 20, ReplicaID: 2, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 3, StoreID: 30, ReplicaID: 3, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 4, StoreID: 40, ReplicaID: 4, Type: roachpb.ReplicaTypeVoterFull()},
		}))

	live := liveness.IsLiveMapEntry{IsLive: true}
	// A dead node under operator decommission: non-active membership, not live.
	deadDecommissioning := liveness.IsLiveMapEntry{
		Liveness: livenesspb.Liveness{Membership: livenesspb.MembershipStatus_DECOMMISSIONING},
		IsLive:   false,
	}

	{
		// The madrid phantom window: 5 healthy voters, all live and
		// membership-active, but clusterNodes reads phantom-low. The floor holds
		// needed at 5, so the range is neither over- nor under-replicated.
		// Without the floor this read over-replicated (needed 3 < live 5).
		ctr, down, under, over := calcRangeCounter(10, fiveVoters, leaseStatus, liveness.IsLiveMap{
			1: live, 2: live, 3: live, 4: live, 5: live,
		}, 5 /* numVoters */, 5 /* numReplicas */, 3 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.False(t, under)
		require.False(t, over)
	}

	{
		// Phantom window with two voters absent from the liveness map (the actual
		// madrid state on n1). Absent voters still count toward the floor, so
		// needed stays 5 against 3 live: under-replicated — the honest signal that
		// the leaseholder cannot account for two of its replicas. Without the
		// floor this read quiet (needed 3 == live 3): blind exactly when it
		// mattered.
		ctr, down, under, over := calcRangeCounter(10, fiveVoters, leaseStatus, liveness.IsLiveMap{
			1: live, 2: live, 3: live,
		}, 5 /* numVoters */, 5 /* numReplicas */, 3 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.True(t, under)
		require.False(t, over)
	}

	{
		// Dead-node decommission in progress: the non-active membership excludes
		// node 5 from the floor (4, matching the allocator's target), so the range
		// reads neither over- nor under-replicated while the dead replica is shed.
		// Without the floor this read over-replicated (needed 3 < live 4) while
		// the allocator was still repairing.
		ctr, down, under, over := calcRangeCounter(10, fiveVoters, leaseStatus, liveness.IsLiveMap{
			1: live, 2: live, 3: live, 4: live, 5: deadDecommissioning,
		}, 5 /* numVoters */, 5 /* numReplicas */, 4 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.False(t, under)
		require.False(t, over)
	}

	{
		// The even steady state: 4 healthy voters on a 4-node cluster. The floor
		// holds needed at 4, so the state the allocator deliberately preserves is
		// not reported as permanently over-replicated (needed 3 < live 4 before).
		ctr, down, under, over := calcRangeCounter(10, fourVoters, leaseStatus, liveness.IsLiveMap{
			1: live, 2: live, 3: live, 4: live,
		}, 5 /* numVoters */, 5 /* numReplicas */, 4 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.False(t, under)
		require.False(t, over)
	}

	{
		// Genuine over-replication versus the configured RF is still detected:
		// the floor is capped at numVoters, so 5 live voters against a configured
		// RF of 3 reads over-replicated.
		ctr, down, under, over := calcRangeCounter(10, fiveVoters, leaseStatus, liveness.IsLiveMap{
			1: live, 2: live, 3: live, 4: live, 5: live,
		}, 3 /* numVoters */, 3 /* numReplicas */, 5 /* clusterNodes */)

		require.True(t, ctr)
		require.False(t, down)
		require.False(t, under)
		require.True(t, over)
	}
}

func TestCalcRangeCounterLeaseHolder(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	rangeDesc := roachpb.NewRangeDescriptor(123, roachpb.RKeyMin, roachpb.RKeyMax,
		roachpb.MakeReplicaSet([]roachpb.ReplicaDescriptor{
			{NodeID: 1, StoreID: 10, ReplicaID: 100, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 2, StoreID: 20, ReplicaID: 200, Type: roachpb.ReplicaTypeNonVoter()},
			{NodeID: 3, StoreID: 30, ReplicaID: 300, Type: roachpb.ReplicaTypeVoterFull()},
			{NodeID: 4, StoreID: 40, ReplicaID: 400, Type: roachpb.ReplicaTypeVoterFull()},
		}))

	leaseStatus := kvserverpb.LeaseStatus{
		Lease: roachpb.Lease{
			Replica: roachpb.ReplicaDescriptor{
				NodeID:    3,
				StoreID:   30,
				ReplicaID: 300,
			},
		},
		State: kvserverpb.LeaseState_VALID,
	}
	leaseStatusInvalid := kvserverpb.LeaseStatus{
		Lease: roachpb.Lease{
			Replica: roachpb.ReplicaDescriptor{
				NodeID:    3,
				StoreID:   30,
				ReplicaID: 300,
			},
		},
		State: kvserverpb.LeaseState_ERROR,
	}

	testcases := []struct {
		desc          string
		storeID       roachpb.StoreID
		leaseStatus   kvserverpb.LeaseStatus
		liveNodes     []roachpb.NodeID
		expectCounter bool
	}{
		{
			desc:          "leaseholder is counter",
			storeID:       30,
			leaseStatus:   leaseStatus,
			liveNodes:     []roachpb.NodeID{1, 2, 3, 4},
			expectCounter: true,
		},
		{
			desc:          "invalid leaseholder is counter",
			storeID:       30,
			leaseStatus:   leaseStatusInvalid,
			liveNodes:     []roachpb.NodeID{1, 2, 3, 4},
			expectCounter: true,
		},
		{
			desc:          "non-leaseholder is not counter",
			storeID:       10,
			leaseStatus:   leaseStatus,
			liveNodes:     []roachpb.NodeID{1, 2, 3, 4},
			expectCounter: false,
		},
		{
			desc:          "non-leaseholder not counter with invalid lease",
			storeID:       10,
			leaseStatus:   leaseStatusInvalid,
			liveNodes:     []roachpb.NodeID{1, 2, 3, 4},
			expectCounter: false,
		},
		{
			desc:          "unavailable leaseholder is not counter",
			storeID:       30,
			leaseStatus:   leaseStatus,
			liveNodes:     []roachpb.NodeID{1, 2, 4},
			expectCounter: false,
		},
		{
			desc:          "first is counter with unavailable leaseholder",
			storeID:       10,
			leaseStatus:   leaseStatus,
			liveNodes:     []roachpb.NodeID{1, 2, 4},
			expectCounter: true,
		},
		{
			desc:          "other is not counter with unavailable leaseholder",
			storeID:       20,
			leaseStatus:   leaseStatus,
			liveNodes:     []roachpb.NodeID{1, 2, 4},
			expectCounter: false,
		},
		{
			desc:          "non-voter can be counter",
			storeID:       20,
			leaseStatus:   leaseStatus,
			liveNodes:     []roachpb.NodeID{2, 4},
			expectCounter: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			livenessMap := liveness.IsLiveMap{}
			for _, nodeID := range tc.liveNodes {
				livenessMap[nodeID] = liveness.IsLiveMapEntry{IsLive: true}
			}
			ctr, _, _, _ := calcRangeCounter(tc.storeID, rangeDesc, tc.leaseStatus, livenessMap,
				3 /* numVoters */, 4 /* numReplicas */, 4 /* clusterNodes */)
			require.Equal(t, tc.expectCounter, ctr)
		})
	}
}
