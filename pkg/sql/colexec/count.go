// Copyright 2018 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// countOp is an operator that counts the number of input rows it receives,
// consuming its entire input and outputting a batch with a single integer
// column containing a single integer, the count of rows received from the
// upstream.
type countOp struct {
	colexecop.OneInputHelper

	internalBatch coldata.Batch
	done          bool
	count         int64
}

var _ colexecop.Operator = &countOp{}

// NewCountOp returns a new count operator that counts the rows in its input.
func NewCountOp(allocator *colmem.Allocator, input colexecop.Operator) colexecop.Operator {
	c := &countOp{
		OneInputHelper: colexecop.MakeOneInputHelper(input),
	}
	c.internalBatch = allocator.NewMemBatchWithFixedCapacity(
		[]*types.T{types.Int}, 1, /* capacity */
	)
	return c
}

func (c *countOp) Next() coldata.Batch {
	if c.done {
		return coldata.ZeroBatch
	}
	c.internalBatch.ResetInternalBatch()
	for {
		bat := c.Input.Next()
		length := bat.Length()
		if length == 0 {
			c.done = true
			c.internalBatch.ColVec(0).Int64()[0] = c.count
			c.internalBatch.SetLength(1)
			return c.internalBatch
		}
		c.count += int64(length)
	}
}
