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
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

type targetSpec struct {
	from, to        scpb.Status
	transitionSpecs []transitionSpec
}

type transitionSpec struct {
	from       scpb.Status
	to         scpb.Status
	revertible bool
	minPhase   scop.Phase
	emitFns    []interface{}
}

type transitionProperty interface {
	apply(spec *transitionSpec)
}

func to(to scpb.Status, properties ...transitionProperty) transitionSpec {
	ts := transitionSpec{
		to:         to,
		revertible: true,
	}
	for _, p := range properties {
		p.apply(&ts)
	}
	return ts
}

func revertible(b bool) transitionProperty {
	return revertibleProperty(b)
}

func minPhase(p scop.Phase) transitionProperty {
	return phaseProperty(p)
}

func emit(fn interface{}) transitionProperty {
	return emitFnSpec{fn}
}

type phaseProperty scop.Phase

func (p phaseProperty) apply(spec *transitionSpec) {
	spec.minPhase = scop.Phase(p)
}

type revertibleProperty bool

func (r revertibleProperty) apply(spec *transitionSpec) {
	spec.revertible = bool(r)
}

var _ transitionProperty = revertibleProperty(true)

type emitFnSpec struct {
	fn interface{}
}

func (e emitFnSpec) apply(spec *transitionSpec) {
	spec.emitFns = append(spec.emitFns, e.fn)
}
