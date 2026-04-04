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
)

// limitOp is an operator that implements limit, returning only the first n
// tuples from its input.
type limitOp struct {
	colexecop.OneInputInitCloserHelper

	limit uint64

	// seen is the number of tuples seen so far.
	seen uint64
	// done is true if the limit has been reached.
	done bool
}

var _ colexecop.Operator = &limitOp{}
var _ colexecop.ClosableOperator = &limitOp{}

// NewLimitOp returns a new limit operator with the given limit.
func NewLimitOp(input colexecop.Operator, limit uint64) colexecop.Operator {
	c := &limitOp{
		OneInputInitCloserHelper: colexecop.MakeOneInputInitCloserHelper(input),
		limit:                    limit,
	}
	return c
}

func (c *limitOp) Next() coldata.Batch {
	if c.done {
		return coldata.ZeroBatch
	}
	bat := c.Input.Next()
	length := bat.Length()
	if length == 0 {
		return bat
	}
	newSeen := c.seen + uint64(length)
	if newSeen >= c.limit {
		c.done = true
		bat.SetLength(int(c.limit - c.seen))
		return bat
	}
	c.seen = newSeen
	return bat
}
