// Copyright 2020 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colconv"
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexecutils"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// NewTupleProjOp creates a new tupleProjOp that projects newly-created tuples
// at position outputIdx taking the tuples' contents from corresponding values
// of the vectors at positions tupleContentsIdxs.
func NewTupleProjOp(
	allocator *colmem.Allocator,
	inputTypes []*types.T,
	tupleContentsIdxs []int,
	outputType *types.T,
	input colexecop.Operator,
	outputIdx int,
) colexecop.Operator {
	input = colexecutils.NewVectorTypeEnforcer(allocator, input, outputType, outputIdx)
	return &tupleProjOp{
		OneInputHelper:    colexecop.MakeOneInputHelper(input),
		allocator:         allocator,
		converter:         colconv.NewVecToDatumConverter(len(inputTypes), tupleContentsIdxs, true /* willRelease */),
		tupleContentsIdxs: tupleContentsIdxs,
		outputType:        outputType,
		outputIdx:         outputIdx,
	}
}

type tupleProjOp struct {
	colexecop.OneInputHelper

	allocator         *colmem.Allocator
	converter         *colconv.VecToDatumConverter
	tupleContentsIdxs []int
	outputType        *types.T
	outputIdx         int
}

var _ colexecop.Operator = &tupleProjOp{}
var _ execinfra.Releasable = &tupleProjOp{}

func (t *tupleProjOp) Next() coldata.Batch {
	batch := t.Input.Next()
	n := batch.Length()
	if n == 0 {
		return coldata.ZeroBatch
	}
	t.converter.ConvertBatchAndDeselect(batch)
	projVec := batch.ColVec(t.outputIdx)
	if projVec.MaybeHasNulls() {
		// We need to make sure that there are no left over null values in the
		// output vector.
		projVec.Nulls().UnsetNulls()
	}

	t.allocator.PerformOperation([]coldata.Vec{projVec}, func() {
		// Preallocate the tuples and their underlying datums in a contiguous
		// slice to reduce allocations.
		tuples := make([]tree.DTuple, n)
		l := len(t.tupleContentsIdxs)
		datums := make(tree.Datums, n*l)
		projCol := projVec.Datum()
		projectInto := func(dst, src int) {
			tuples[src] = tree.MakeDTuple(
				t.outputType, datums[src*l:(src+1)*l:(src+1)*l]...,
			)
			projCol.Set(dst, t.projectInto(&tuples[src], src))
		}
		if sel := batch.Selection(); sel != nil {
			for convertedIdx, i := range sel[:n] {
				projectInto(i, convertedIdx)
			}
		} else {
			for i := 0; i < n; i++ {
				projectInto(i, i)
			}
		}
	})
	return batch
}

func (t *tupleProjOp) projectInto(tuple *tree.DTuple, convertedIdx int) tree.Datum {
	for i, columnIdx := range t.tupleContentsIdxs {
		tuple.D[i] = t.converter.GetDatumColumn(columnIdx)[convertedIdx]
	}
	return tuple
}

// Release is part of the execinfra.Releasable interface.
func (t *tupleProjOp) Release() {
	t.converter.Release()
}
