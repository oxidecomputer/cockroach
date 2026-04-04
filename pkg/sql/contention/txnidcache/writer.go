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
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/contentionpb"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
)

// There is no strong reason why shardCount is 16 beyond that Java's
// ConcurrentHashMap also uses 16 shards and has reasonably good performance.
const shardCount = 16

type writer struct {
	st *cluster.Settings

	shards [shardCount]*concurrentWriteBuffer

	sink blockSink
}

var _ Writer = &writer{}

func newWriter(st *cluster.Settings, sink blockSink) *writer {
	w := &writer{
		st:   st,
		sink: sink,
	}

	for shardIdx := 0; shardIdx < shardCount; shardIdx++ {
		w.shards[shardIdx] = newConcurrentWriteBuffer(sink)
	}

	return w
}

// Record implements the Writer interface.
func (w *writer) Record(resolvedTxnID contentionpb.ResolvedTxnID) {
	if MaxSize.Get(&w.st.SV) == 0 {
		return
	}

	// There are edge cases where the txnID in the resolvedTxnID will be nil,
	// (e.g. when the connExecutor closes while a transaction is still active).
	// This causes that, occasionally, connExecutor will emit resolvedTxnIDs with
	// invalid txnID but valid txnFingerprintID. Writing invalid txnID into the
	// writer can potentially cause data loss. (Since the TxnID cache stops
	// processing the input batch when it encounters the first invalid txnID).
	if resolvedTxnID.TxnID.Equal(uuid.Nil) {
		return
	}
	shardIdx := hashTxnID(resolvedTxnID.TxnID)
	buffer := w.shards[shardIdx]
	buffer.Record(resolvedTxnID)
}

// DrainWriteBuffer implements the Writer interface.
func (w *writer) DrainWriteBuffer() {
	for shardIdx := 0; shardIdx < shardCount; shardIdx++ {
		w.shards[shardIdx].DrainWriteBuffer()
	}
}

func hashTxnID(txnID uuid.UUID) int {
	b := txnID.GetBytes()
	_, val, err := encoding.DecodeUint64Descending(b)
	if err != nil {
		panic(err)
	}
	return int(val % shardCount)
}
