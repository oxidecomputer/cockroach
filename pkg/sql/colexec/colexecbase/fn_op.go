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

package colexecbase

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
)

// fnOp is an operator that executes an arbitrary function for its side-effects,
// once per input batch, passing the input batch unmodified along.
type fnOp struct {
	colexecop.OneInputHelper
	colexecop.NonExplainable

	fn func()
}

var _ colexecop.ResettableOperator = &fnOp{}

func (f *fnOp) Next() coldata.Batch {
	batch := f.Input.Next()
	f.fn()
	return batch
}

func (f *fnOp) Reset(ctx context.Context) {
	if resettableOp, ok := f.Input.(colexecop.Resetter); ok {
		resettableOp.Reset(ctx)
	}
}
