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

package colexec

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
	"github.com/cockroachdb/cockroach/pkg/util/buildutil"
	"github.com/cockroachdb/errors"
)

// invariantsChecker is a helper Operator that will check that invariants that
// are present in the vectorized engine are maintained on all batches. It
// should be planned between other Operators in tests.
type invariantsChecker struct {
	colexecop.OneInputNode
	colexecop.InitHelper
	colexecop.NonExplainable

	metadataSource colexecop.MetadataSource
}

var _ colexecop.DrainableOperator = &invariantsChecker{}
var _ colexecop.ClosableOperator = &invariantsChecker{}

// NewInvariantsChecker creates a new invariantsChecker.
func NewInvariantsChecker(input colexecop.Operator) colexecop.DrainableOperator {
	if !buildutil.CrdbTestBuild {
		colexecerror.InternalError(errors.AssertionFailedf(
			"an invariantsChecker is attempted to be created in non-test build",
		))
	}
	c := &invariantsChecker{
		OneInputNode: colexecop.OneInputNode{Input: input},
	}
	if ms, ok := input.(colexecop.MetadataSource); ok {
		c.metadataSource = ms
	}
	return c
}

// MaybeUnwrapInvariantsChecker checks whether op is an invariants checker and
// returns its input if so, otherwise op is returned.
func MaybeUnwrapInvariantsChecker(op colexecop.Operator) colexecop.Operator {
	if i, ok := op.(*invariantsChecker); ok {
		return i.Input
	}
	return op
}

// Init implements the colexecop.Operator interface.
func (i *invariantsChecker) Init(ctx context.Context) {
	if !i.InitHelper.Init(ctx) {
		return
	}
	i.Input.Init(i.Ctx)
}

// assertInitWasCalled asserts that Init() has been called on the invariants
// checker and returns a boolean indicating whether the execution should be
// short-circuited (true means that the caller should just return right away).
func (i *invariantsChecker) assertInitWasCalled() bool {
	if i.Ctx == nil {
		if c, ok := i.Input.(*Columnarizer); ok {
			if c.removedFromFlow {
				// This is a special case in which we allow for the operator to
				// not be initialized. Next and DrainMeta calls are noops in
				// this case, so the caller should short-circuit.
				return true
			}
		}
		colexecerror.InternalError(errors.AssertionFailedf("Init hasn't been called, input is %T", i.Input))
	}
	return false
}

// Next implements the colexecop.Operator interface.
func (i *invariantsChecker) Next() coldata.Batch {
	if shortCircuit := i.assertInitWasCalled(); shortCircuit {
		return coldata.ZeroBatch
	}
	b := i.Input.Next()
	n := b.Length()
	if n == 0 {
		return b
	}
	if sel := b.Selection(); sel != nil {
		for i := 1; i < n; i++ {
			if sel[i] <= sel[i-1] {
				colexecerror.InternalError(errors.AssertionFailedf(
					"unexpectedly selection vector is not an increasing sequence "+
						"at position %d: %v", i, sel[:n],
				))
			}
		}
	}
	return b
}

// DrainMeta implements the colexecop.MetadataSource interface.
func (i *invariantsChecker) DrainMeta() []execinfrapb.ProducerMetadata {
	if shortCircuit := i.assertInitWasCalled(); shortCircuit {
		return nil
	}
	if i.metadataSource == nil {
		return nil
	}
	return i.metadataSource.DrainMeta()
}

// Close is part of the colexecop.ClosableOperator interface.
func (i *invariantsChecker) Close(ctx context.Context) error {
	c, ok := i.Input.(colexecop.Closer)
	if !ok {
		return nil
	}
	return c.Close(ctx)
}
