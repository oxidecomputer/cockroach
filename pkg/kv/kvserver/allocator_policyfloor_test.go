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
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness/livenesspb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// TODO-RAINCLAUDE: experiment 2 for omicron#10658. Verifies the policy-floor:
// the effective RF is driven only by operator policy (configured RF +
// decommissioning), never by a transient/cold-cache node count. Contrast with
// experiment 1 (drop the downshift): the floor keeps the downshift's *good*
// behavior — a genuinely small cluster stays happy instead of churning — while
// still neutralizing the madrid trim.
func TestAllocatorPolicyFloorScenarios(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	var numNodes int
	stopper, _, _, sp, _ := createTestStorePool(ctx,
		TestTimeUntilStoreDeadOff, false, /* deterministic */
		func() int { return numNodes },
		livenesspb.NodeLivenessStatus_LIVE)
	defer stopper.Stop(ctx)
	a := MakeAllocator(sp, func(string) (time.Duration, bool) {
		return 0, true
	}, nil /* knobs */, nil /* storeMetrics */)
	conf := roachpb.SpanConfig{NumReplicas: 5}

	type scenario struct {
		name            string
		storeList       []roachpb.StoreID
		live            []roachpb.StoreID
		unavailable     []roachpb.StoreID
		dead            []roachpb.StoreID
		decommissioning []roachpb.StoreID
		nodeCount       int
		wantAction      AllocatorAction
		// what experiment 1 (bare drop of the downshift) produced, for contrast
		exp1Action AllocatorAction
		// what the original downshifting code produced
		oldAction AllocatorAction
	}

	scenarios := []scenario{
		{
			name:       "madrid cold cache: 5 live healthy, 0 decommissioning, phantom count 3",
			storeList:  []roachpb.StoreID{1, 2, 3, 4, 5},
			live:       []roachpb.StoreID{1, 2, 3, 4, 5},
			nodeCount:  3,
			wantAction: AllocatorConsiderRebalance, // FIXED (floored to 5)
			exp1Action: AllocatorConsiderRebalance,
			oldAction:  AllocatorRemoveVoter, // the bug
		},
		{
			name:            "decommission: range has 5 voters, 2 decommissioning, count 3",
			storeList:       []roachpb.StoreID{1, 2, 3, 4, 5},
			live:            []roachpb.StoreID{1, 2, 3},
			decommissioning: []roachpb.StoreID{4, 5},
			nodeCount:       3,
			wantAction:      AllocatorRemoveDecommissioningVoter, // proceeds
			exp1Action:      AllocatorReplaceDecommissioningVoter,
			oldAction:       AllocatorRemoveDecommissioningVoter,
		},
		{
			name:       "small cluster: RF5 range on 3 nodes, only 3 exist, 0 decommissioning",
			storeList:  []roachpb.StoreID{1, 2, 3},
			live:       []roachpb.StoreID{1, 2, 3},
			nodeCount:  3,
			wantAction: AllocatorConsiderRebalance, // happy — NO churn (floor==downshift==3)
			exp1Action: AllocatorAddVoter,          // experiment 1 churned here
			oldAction:  AllocatorConsiderRebalance,
		},
		{
			name:       "dead-but-present (control): 2 dead, warm cache counts 5",
			storeList:  []roachpb.StoreID{1, 2, 3, 4, 5},
			live:       []roachpb.StoreID{1, 2, 3},
			dead:       []roachpb.StoreID{4, 5},
			nodeCount:  5,
			wantAction: AllocatorReplaceDeadVoter,
			exp1Action: AllocatorReplaceDeadVoter,
			oldAction:  AllocatorReplaceDeadVoter,
		},
		{
			name:       "dead + cold cache: 2 dead (NOT decommissioned), phantom count 3",
			storeList:  []roachpb.StoreID{1, 2, 3, 4, 5},
			live:       []roachpb.StoreID{1, 2, 3},
			dead:       []roachpb.StoreID{4, 5},
			nodeCount:  3,
			wantAction: AllocatorReplaceDeadVoter, // floored to 5: keep trying to maintain RF, never trim
			exp1Action: AllocatorReplaceDeadVoter,
			oldAction:  AllocatorRemoveDeadVoter, // old code trimmed toward 3
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			numNodes = s.nodeCount
			mockStorePool(sp, s.live, s.unavailable, s.dead, s.decommissioning, nil, nil)
			desc := makeDescriptor(s.storeList)
			desc.EndKey = roachpb.RKey(keys.SystemPrefix)
			action, _ := a.ComputeAction(ctx, conf, &desc)
			require.Equalf(t, s.wantAction.String(), action.String(),
				"policy-floor wants %s (experiment 1 gave %s; original downshift gave %s)",
				s.wantAction, s.exp1Action, s.oldAction)
		})
	}
}

// TODO-RAINCLAUDE: experiment 2 — a full decommission of two nodes from a 5-node
// RF-5 cluster must still complete under the policy floor. Each step removes a
// decommissioning voter until the range settles at 3 and is happy.
func TestAllocatorPolicyFloorDecommissionCompletes(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	var numNodes int
	stopper, _, _, sp, _ := createTestStorePool(ctx,
		TestTimeUntilStoreDeadOff, false, /* deterministic */
		func() int { return numNodes },
		livenesspb.NodeLivenessStatus_LIVE)
	defer stopper.Stop(ctx)
	a := MakeAllocator(sp, func(string) (time.Duration, bool) {
		return 0, true
	}, nil /* knobs */, nil /* storeMetrics */)
	conf := roachpb.SpanConfig{NumReplicas: 5}

	// Operator decommissions nodes 4 and 5 of a 5-node cluster. 3 nodes remain
	// (count excludes decommissioning), so the cluster genuinely shrinks below RF.
	numNodes = 3

	steps := []struct {
		storeList       []roachpb.StoreID
		decommissioning []roachpb.StoreID
		live            []roachpb.StoreID
		wantAction      AllocatorAction
	}{
		{[]roachpb.StoreID{1, 2, 3, 4, 5}, []roachpb.StoreID{4, 5}, []roachpb.StoreID{1, 2, 3}, AllocatorRemoveDecommissioningVoter},
		{[]roachpb.StoreID{1, 2, 3, 5}, []roachpb.StoreID{5}, []roachpb.StoreID{1, 2, 3}, AllocatorRemoveDecommissioningVoter},
		{[]roachpb.StoreID{1, 2, 3}, nil, []roachpb.StoreID{1, 2, 3}, AllocatorConsiderRebalance}, // settled
	}

	for i, s := range steps {
		mockStorePool(sp, s.live, nil, nil, s.decommissioning, nil, nil)
		desc := makeDescriptor(s.storeList)
		desc.EndKey = roachpb.RKey(keys.SystemPrefix)
		action, _ := a.ComputeAction(ctx, conf, &desc)
		require.Equalf(t, s.wantAction.String(), action.String(), "decommission step %d", i)
	}
}
