// Copyright 2019 The Cockroach Authors.
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

package constraint

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// AnalyzedConstraints represents the result or AnalyzeConstraints(). It
// combines a zone's constraints with information about which stores satisfy
// what term of the constraints disjunction.
type AnalyzedConstraints struct {
	Constraints []roachpb.ConstraintsConjunction
	// True if the per-replica constraints don't fully cover all the desired
	// replicas in the range (sum(constraints.NumReplicas) < zone.NumReplicas).
	// In such cases, we allow replicas that don't match any of the per-replica
	// constraints, but never mark them as necessary.
	UnconstrainedReplicas bool
	// For each conjunction of constraints in the above slice, track which
	// StoreIDs satisfy them. This field is unused if there are no constraints.
	SatisfiedBy [][]roachpb.StoreID
	// Maps from StoreID to the indices in the constraints slice of which
	// constraints the store satisfies. This field is unused if there are no
	// constraints.
	Satisfies map[roachpb.StoreID][]int
}

// EmptyAnalyzedConstraints represents an empty set of constraints that are
// satisfied by any given configuration of replicas.
var EmptyAnalyzedConstraints = AnalyzedConstraints{}

// AnalyzeConstraints processes the zone config constraints that apply to a
// range along with the current replicas for a range, spitting back out
// information about which constraints are satisfied by which replicas and
// which replicas satisfy which constraints, aiding in allocation decisions.
func AnalyzeConstraints(
	ctx context.Context,
	getStoreDescFn func(roachpb.StoreID) (roachpb.StoreDescriptor, bool),
	existing []roachpb.ReplicaDescriptor,
	numReplicas int32,
	constraints []roachpb.ConstraintsConjunction,
) AnalyzedConstraints {
	result := AnalyzedConstraints{
		Constraints: constraints,
	}

	if len(constraints) > 0 {
		result.SatisfiedBy = make([][]roachpb.StoreID, len(constraints))
		result.Satisfies = make(map[roachpb.StoreID][]int)
	}

	var constrainedReplicas int32
	for i, subConstraints := range constraints {
		constrainedReplicas += subConstraints.NumReplicas
		for _, repl := range existing {
			// If for some reason we don't have the store descriptor (which shouldn't
			// happen once a node is hooked into gossip), trust that it's valid. This
			// is a much more stable failure state than frantically moving everything
			// off such a node.
			store, ok := getStoreDescFn(repl.StoreID)
			if !ok || ConjunctionsCheck(store, subConstraints.Constraints) {
				result.SatisfiedBy[i] = append(result.SatisfiedBy[i], store.StoreID)
				result.Satisfies[store.StoreID] = append(result.Satisfies[store.StoreID], i)
			}
		}
	}
	if constrainedReplicas > 0 && constrainedReplicas < numReplicas {
		result.UnconstrainedReplicas = true
	}
	return result
}

// ConjunctionsCheck checks a store against a single set of constraints (out of
// the possibly numerous sets that apply to a range), returning true iff the
// store matches the constraints. The contraints are AND'ed together; a store
// matches the conjunction if it matches all of them.
func ConjunctionsCheck(store roachpb.StoreDescriptor, constraints []roachpb.Constraint) bool {
	for _, constraint := range constraints {
		// StoreMatchesConstraint returns whether a store matches the given constraint.
		hasConstraint := roachpb.StoreMatchesConstraint(store, constraint)
		if (constraint.Type == roachpb.Constraint_REQUIRED && !hasConstraint) ||
			(constraint.Type == roachpb.Constraint_PROHIBITED && hasConstraint) {
			return false
		}
	}
	return true
}
