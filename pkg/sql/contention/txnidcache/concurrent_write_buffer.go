// Copyright 2022 The Cockroach Authors.
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

package txnidcache

import (
	"sync"

	"github.com/cockroachdb/cockroach/pkg/sql/contention/contentionutils"
	"github.com/cockroachdb/cockroach/pkg/sql/contentionpb"
)

// blockSize is chosen as 168 since each ResolvedTxnID is 24 byte.
// 168 * 24 = 4032 bytes < 4KiB page size.
const blockSize = 168

type block [blockSize]contentionpb.ResolvedTxnID

var blockPool = &sync.Pool{
	New: func() interface{} {
		return &block{}
	},
}

// concurrentWriteBuffer is a data structure that optimizes for concurrent
// writes and also implements the Writer interface.
type concurrentWriteBuffer struct {
	guard struct {
		*contentionutils.ConcurrentBufferGuard

		// block is the temporary buffer that concurrentWriteBuffer uses to batch
		// write requests before sending them into the channel.
		block *block
	}

	// sink is the flush target that ConcurrentWriteBuffer flushes to once
	// block is full.
	sink blockSink
}

var _ Writer = &concurrentWriteBuffer{}

// newConcurrentWriteBuffer returns a new instance of concurrentWriteBuffer.
func newConcurrentWriteBuffer(sink blockSink) *concurrentWriteBuffer {
	writeBuffer := &concurrentWriteBuffer{
		sink: sink,
	}

	writeBuffer.guard.block = blockPool.Get().(*block)
	writeBuffer.guard.ConcurrentBufferGuard = contentionutils.NewConcurrentBufferGuard(
		func() int64 {
			return blockSize
		}, /* limiter */
		func(_ int64) {
			writeBuffer.sink.push(writeBuffer.guard.block)

			// Resets the block.
			writeBuffer.guard.block = blockPool.Get().(*block)
		} /* onBufferFull */)

	return writeBuffer
}

// Record records a mapping from txnID to its corresponding transaction
// fingerprint ID. Record is safe to be used concurrently.
func (c *concurrentWriteBuffer) Record(resolvedTxnID contentionpb.ResolvedTxnID) {
	c.guard.AtomicWrite(func(writerIdx int64) {
		c.guard.block[writerIdx] = resolvedTxnID
	})
}

// DrainWriteBuffer flushes concurrentWriteBuffer into the channel. It
// implements the txnidcache.Writer interface.
func (c *concurrentWriteBuffer) DrainWriteBuffer() {
	c.guard.ForceSync()
}
