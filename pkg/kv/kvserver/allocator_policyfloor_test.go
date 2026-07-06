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
// still neutralizing the madrid trim. The floor is keyed on liveness membership
// (DECOMMISSIONING or DECOMMISSIONED), not store-pool status, so that a DEAD
// node under decommission (which reports DECOMMISSIONED -> storeStatusDead)
// still lowers it; see the dead-node decommission scenarios.
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
		// decommissioned mocks NodeLivenessStatus_DECOMMISSIONED, which is what a
		// DEAD node with DECOMMISSIONING (or DECOMMISSIONED) membership reports —
		// it classifies as storeStatusDead, not storeStatusDecommissioning.
		decommissioned []roachpb.StoreID
		nodeCount      int
		wantAction     AllocatorAction
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
		{
			// A node died permanently and the operator decommissions it. The
			// membership-keyed floor excludes it (floor 4, not 5), so the dead
			// replica is REMOVED and the decommission completes even with no spare
			// node. The status-keyed floor (the original experiment 2 formula) kept
			// needed at 5 and produced ReplaceDeadVoter, which has no allocation
			// target on the 4-node remainder: purgatory forever, and the
			// decommission never reaches zero replicas.
			name:           "dead-node decommission, no spare: 1 decommissioned(dead), count 4",
			storeList:      []roachpb.StoreID{1, 2, 3, 4, 5},
			live:           []roachpb.StoreID{1, 2, 3, 4},
			decommissioned: []roachpb.StoreID{5},
			nodeCount:      4,
			wantAction:     AllocatorRemoveDeadVoter,
			exp1Action:     AllocatorReplaceDeadVoter, // stalls: no 5th node
			oldAction:      AllocatorRemoveDeadVoter,
		},
		{
			// Same, but a spare node 6 exists (the Oxide add-then-decommission
			// flow). The warm count keeps needed at 5 above the floor of 4, so the
			// dead replica is replaced onto the spare and the range stays at the
			// configured RF.
			name:           "dead-node decommission, spare exists: count 5",
			storeList:      []roachpb.StoreID{1, 2, 3, 4, 5},
			live:           []roachpb.StoreID{1, 2, 3, 4, 6},
			decommissioned: []roachpb.StoreID{5},
			nodeCount:      5,
			wantAction:     AllocatorReplaceDeadVoter,
			exp1Action:     AllocatorReplaceDeadVoter,
			oldAction:      AllocatorReplaceDeadVoter,
		},
		{
			// A decommissioning node that is transiently unreachable reports
			// UNAVAILABLE (LivenessStatus collapses membership for expired-but-not-
			// yet-dead nodes) -> storeStatusUnknown. It is deliberately NOT excluded
			// from the floor: the decommission pauses until the node resolves to
			// live (remove decommissioning) or dead (remove dead). Excluding it
			// would open a liveness-blind RemoveVoter window. Note the old code
			// trimmed a LIVE voter here — the madrid mechanism without any cache
			// coldness on the affected node.
			name:        "decommissioning node transiently unavailable, phantom count 4",
			storeList:   []roachpb.StoreID{1, 2, 3, 4, 5},
			live:        []roachpb.StoreID{1, 2, 3, 4},
			unavailable: []roachpb.StoreID{5},
			nodeCount:   4,
			wantAction:  AllocatorConsiderRebalance, // pause, fail-safe
			exp1Action:  AllocatorConsiderRebalance,
			oldAction:   AllocatorRemoveVoter, // blind trim of a live voter
		},
		{
			// The even steady state after a 5->4 shrink: the floor holds 4 healthy
			// voters on a 4-node cluster rather than trimming to 3. The even-quorum
			// nicety yields to the don't-trim-healthy-replicas invariant; the
			// operator escape hatch is lowering num_replicas (policy, honored via
			// the cap).
			name:       "even steady state: 4 healthy voters, 4 nodes",
			storeList:  []roachpb.StoreID{1, 2, 3, 4},
			live:       []roachpb.StoreID{1, 2, 3, 4},
			nodeCount:  4,
			wantAction: AllocatorConsiderRebalance,
			exp1Action: AllocatorAddVoter,    // wants a 5th node that does not exist
			oldAction:  AllocatorRemoveVoter, // trimmed the even 4 to 3
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			numNodes = s.nodeCount
			mockStorePool(sp, s.live, s.unavailable, s.dead, s.decommissioning, s.decommissioned, nil)
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

// TODO-RAINCLAUDE: membership-keyed floor fix — decommissioning a DEAD node
// must complete. Without a spare node the dead replica is removed and the range
// settles at 4 voters (the even state the floor preserves); with a spare it is
// replaced and the range stays at the configured RF 5. Before the fix the
// no-spare path returned ReplaceDeadVoter forever (no allocation target ->
// purgatory) and the node could never finish decommissioning.
func TestAllocatorPolicyFloorDeadNodeDecommission(t *testing.T) {
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

	steps := []struct {
		name           string
		storeList      []roachpb.StoreID
		live           []roachpb.StoreID
		decommissioned []roachpb.StoreID
		nodeCount      int
		wantAction     AllocatorAction
	}{
		// No spare: node 5 is dead and being decommissioned on what is now a
		// 4-node cluster. The dead replica is removed, then the range settles at
		// the even 4-voter state.
		{
			name:           "no spare: shed the dead replica",
			storeList:      []roachpb.StoreID{1, 2, 3, 4, 5},
			live:           []roachpb.StoreID{1, 2, 3, 4},
			decommissioned: []roachpb.StoreID{5},
			nodeCount:      4,
			wantAction:     AllocatorRemoveDeadVoter,
		},
		{
			name:           "no spare: settled at 4",
			storeList:      []roachpb.StoreID{1, 2, 3, 4},
			live:           []roachpb.StoreID{1, 2, 3, 4},
			decommissioned: []roachpb.StoreID{5},
			nodeCount:      4,
			wantAction:     AllocatorConsiderRebalance,
		},
		// With a spare (the Oxide add-then-decommission flow): the dead replica is
		// replaced onto node 6, then the range is settled at the configured RF.
		{
			name:           "spare: replace the dead replica",
			storeList:      []roachpb.StoreID{1, 2, 3, 4, 5},
			live:           []roachpb.StoreID{1, 2, 3, 4, 6},
			decommissioned: []roachpb.StoreID{5},
			nodeCount:      5,
			wantAction:     AllocatorReplaceDeadVoter,
		},
		{
			name:           "spare: settled at RF 5",
			storeList:      []roachpb.StoreID{1, 2, 3, 4, 6},
			live:           []roachpb.StoreID{1, 2, 3, 4, 6},
			decommissioned: []roachpb.StoreID{5},
			nodeCount:      5,
			wantAction:     AllocatorConsiderRebalance,
		},
	}

	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			numNodes = s.nodeCount
			mockStorePool(sp, s.live, nil, nil, nil, s.decommissioned, nil)
			desc := makeDescriptor(s.storeList)
			desc.EndKey = roachpb.RKey(keys.SystemPrefix)
			action, _ := a.ComputeAction(ctx, conf, &desc)
			require.Equalf(t, s.wantAction.String(), action.String(), "step %q", s.name)
		})
	}
}
