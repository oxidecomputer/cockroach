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
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/liveness/livenesspb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// TODO-RAINCLAUDE: omicron#10658 — the policy floor, in one place. The rule:
// the effective replication factor is a function of operator policy alone
// (the configured num_replicas and operator decommissioning), never of
// transient cluster health (a dead node, a cold liveness cache, a phantom-low
// node count). computeAction applies the floor to the allocator's target and
// calcRangeCounter applies it to the under-/over-replication gauges, both
// through the pure helpers below so the two sites cannot drift and the logic
// is directly unit-testable (see the TestPolicyFloor* tests).

// TODO-RAINCLAUDE: the floor arithmetic. Returns neededVoters floored at the
// range's own voters on nodes the operator has not directed out of the
// cluster (haveVoters - policyRemovedVoters), capped at the configured RF.
// Pure; TestPolicyFloorNeededVotersInvariants proves the following properties
// by exhaustive enumeration:
//
//   - The result never drops below neededVoters, so the downshift's
//     small-cluster and decommission-completion behavior is preserved.
//   - The result never exceeds max(neededVoters, configuredVoters), so
//     lowering num_replicas still down-replicates.
//   - The result never exceeds max(neededVoters, haveVoters -
//     policyRemovedVoters): the floor cannot invent replicas.
//   - With no policy-removed voters and haveVoters within the configured RF,
//     the result is at least haveVoters — the madrid invariant: no phantom-low
//     node count can size a healthy range below the replicas it already has,
//     which is what authorized the omicron#10658 trim.
func policyFloorNeededVoters(
	neededVoters, haveVoters, policyRemovedVoters, configuredVoters int,
) int {
	policyFloor := haveVoters - policyRemovedVoters
	if policyFloor > configuredVoters {
		policyFloor = configuredVoters
	}
	if policyFloor > neededVoters {
		return policyFloor
	}
	return neededVoters
}

// TODO-RAINCLAUDE: the allocator-side exclusion predicate: is this node's
// liveness status the result of operator policy (decommissioning), such that
// its replica should not count toward the policy floor? True for exactly
// DECOMMISSIONING (live + non-active membership) and DECOMMISSIONED (dead +
// non-active membership; this is what a dead node under decommission reports,
// and it classifies as storeStatusDead — which is why the floor cannot key on
// storeStatusDecommissioning). Everything else is false, and false is always
// the fail-safe direction: an unexcluded replica keeps the floor high, which
// can pause a decommission but can never remove a live replica. In
// particular: UNKNOWN (no liveness record — the cold-cache state that
// triggered omicron#10658) and UNAVAILABLE (expired but not yet past the dead
// threshold; LivenessStatus reports this even for a decommissioning node, and
// excluding it here would open a liveness-blind RemoveVoter window in
// computeAction while the node sits in neither the dead nor the
// decommissioning status set). The default arm makes any future status
// fail-safe too; TestNodeLivenessStatusIsPolicyRemoved forces an explicit
// classification when a status is added.
func nodeLivenessStatusIsPolicyRemoved(status livenesspb.NodeLivenessStatus) bool {
	switch status {
	case livenesspb.NodeLivenessStatus_DECOMMISSIONING,
		livenesspb.NodeLivenessStatus_DECOMMISSIONED:
		return true
	case livenesspb.NodeLivenessStatus_UNKNOWN,
		livenesspb.NodeLivenessStatus_DEAD,
		livenesspb.NodeLivenessStatus_UNAVAILABLE,
		livenesspb.NodeLivenessStatus_LIVE,
		livenesspb.NodeLivenessStatus_DRAINING:
		return false
	default:
		return false
	}
}

// TODO-RAINCLAUDE: the metric-side exclusion count for calcRangeCounter,
// which has an IsLiveMap rather than a NodeLivenessFunc. Keyed on raw
// liveness membership because the map carries no dead-threshold
// classification; this diverges from nodeLivenessStatusIsPolicyRemoved only
// while a decommissioning node is unavailable-but-not-yet-dead (the allocator
// pauses; the gauge already treats the node as policy-removed), which is
// acceptable for an estimated gauge. A voter absent from the map does not
// count as removed — an unverifiable replica must not lower the target — so
// during a cold-cache window a healthy range reads as under-replicated (the
// honest signal that the leaseholder cannot account for its replicas) rather
// than over-replicated.
func policyRemovedVoterCount(
	voters []roachpb.ReplicaDescriptor, livenessMap liveness.IsLiveMap,
) int {
	removed := 0
	for _, rd := range voters {
		if entry, ok := livenessMap[rd.NodeID]; ok && !entry.Membership.Active() {
			removed++
		}
	}
	return removed
}
