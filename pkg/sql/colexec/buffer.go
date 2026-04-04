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
	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
)

// bufferOp is an operator that buffers a single batch at a time from an input,
// and makes it available to be read multiple times by downstream consumers.
type bufferOp struct {
	colexecop.OneInputHelper

	// read is true if someone has read the current batch already.
	read  bool
	batch coldata.Batch
}

var _ colexecop.Operator = &bufferOp{}

// NewBufferOp returns a new bufferOp, initialized to buffer batches from the
// supplied input.
func NewBufferOp(input colexecop.Operator) colexecop.Operator {
	return &bufferOp{
		OneInputHelper: colexecop.MakeOneInputHelper(input),
	}
}

// rewind resets this buffer to be readable again.
// NOTE: it is the caller responsibility to restore the batch into the desired
// state.
func (b *bufferOp) rewind() {
	b.read = false
}

// advance reads the next batch from the input into the buffer, preparing itself
// for reads.
func (b *bufferOp) advance() {
	b.batch = b.Input.Next()
	b.rewind()
}

func (b *bufferOp) Next() coldata.Batch {
	if b.read {
		return coldata.ZeroBatch
	}
	b.read = true
	return b.batch
}
