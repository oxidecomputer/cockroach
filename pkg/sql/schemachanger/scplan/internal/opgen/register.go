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

package opgen

import (
	"reflect"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/errors"
)

// equiv defines the from status as being equivalent to the current status.
func equiv(from scpb.Status) transitionSpec {
	return transitionSpec{from: from, revertible: true}
}

func notImplemented(e scpb.Element) *scop.NotImplemented {
	return &scop.NotImplemented{
		ElementType: reflect.ValueOf(e).Type().Elem().String(),
	}
}

func toPublic(initialStatus scpb.Status, specs ...transitionSpec) targetSpec {
	return asTargetSpec(scpb.Status_PUBLIC, initialStatus, specs...)
}

func toAbsent(initialStatus scpb.Status, specs ...transitionSpec) targetSpec {
	return asTargetSpec(scpb.Status_ABSENT, initialStatus, specs...)
}

func asTargetSpec(to, from scpb.Status, specs ...transitionSpec) targetSpec {
	return targetSpec{from: from, to: to, transitionSpecs: specs}
}

// register constructs all operations edges for a given element.
// Intended to be called during init, register panics on any error.
func (r *registry) register(e scpb.Element, targetSpecs ...targetSpec) {
	onErrPanic := func(err error) {
		if err != nil {
			panic(errors.NewAssertionErrorWithWrappedErrf(err, "element %T", e))
		}
	}
	targets := make([]target, len(targetSpecs))
	for i, spec := range targetSpecs {
		var err error
		targets[i], err = makeTarget(e, spec)
		onErrPanic(err)
	}
	onErrPanic(validateTargets(targets))
	r.targets = append(r.targets, targets...)
}

func validateTargets(targets []target) error {
	allStatuses := map[scpb.Status]bool{}
	targetStatuses := make([]map[scpb.Status]bool, len(targets))
	for i, tgt := range targets {
		m := map[scpb.Status]bool{}
		for _, t := range tgt.transitions {
			m[t.from] = true
			allStatuses[t.from] = true
			m[t.to] = true
			allStatuses[t.to] = true
		}
		targetStatuses[i] = m
	}

	for i, tgt := range targets {
		m := targetStatuses[i]
		for s := range allStatuses {
			if !m[s] {
				return errors.Errorf("target %s: status %s is missing here but is featured in other targets",
					tgt.status.String(), s.String())
			}
		}
	}
	return nil
}
