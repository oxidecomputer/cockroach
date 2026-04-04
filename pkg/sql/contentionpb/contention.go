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

package contentionpb

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
)

const singleIndentation = "  "
const doubleIndentation = singleIndentation + singleIndentation
const tripleIndentation = doubleIndentation + singleIndentation

const contentionEventsStr = "num contention events:"
const cumulativeContentionTimeStr = "cumulative contention time:"
const contendingTxnsStr = "contending txns:"

func (ice IndexContentionEvents) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("tableID=%d indexID=%d\n", ice.TableID, ice.IndexID))
	b.WriteString(fmt.Sprintf("%s%s %d\n", singleIndentation, contentionEventsStr, ice.NumContentionEvents))
	b.WriteString(fmt.Sprintf("%s%s %s\n", singleIndentation, cumulativeContentionTimeStr, ice.CumulativeContentionTime))
	b.WriteString(fmt.Sprintf("%skeys:\n", singleIndentation))
	for i := range ice.Events {
		b.WriteString(ice.Events[i].String())
	}
	return b.String()
}

func (skc SingleKeyContention) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s%s %s\n", doubleIndentation, skc.Key, contendingTxnsStr))
	for i := range skc.Txns {
		b.WriteString(skc.Txns[i].String())
	}
	return b.String()
}

func toString(stx SingleTxnContention, indentation string) string {
	return fmt.Sprintf("%sid=%s count=%d\n", indentation, stx.TxnID, stx.Count)
}

func (stx SingleTxnContention) String() string {
	return toString(stx, tripleIndentation)
}

func (skc SingleNonSQLKeyContention) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("non-SQL key %s %s\n", skc.Key, contendingTxnsStr))
	b.WriteString(fmt.Sprintf("%s%s %d\n", singleIndentation, contentionEventsStr, skc.NumContentionEvents))
	b.WriteString(fmt.Sprintf("%s%s %s\n", singleIndentation, cumulativeContentionTimeStr, skc.CumulativeContentionTime))
	for i := range skc.Txns {
		b.WriteString(toString(skc.Txns[i], doubleIndentation))
	}
	return b.String()
}

// Valid returns if the ResolvedTxnID is valid.
func (r *ResolvedTxnID) Valid() bool {
	return !uuid.Nil.Equal(r.TxnID)
}

// Valid returns if the ExtendedContentionEvent is valid.
func (e *ExtendedContentionEvent) Valid() bool {
	return !uuid.Nil.Equal(e.BlockingEvent.TxnMeta.ID)
}

// Hash returns a hash that's unique to ExtendedContentionEvent using
// blocking txn's txnID, waiting txn's txnID and the event collection timestamp.
func (e *ExtendedContentionEvent) Hash() uint64 {
	hash := util.MakeFNV64()
	hashUUID(e.BlockingEvent.TxnMeta.ID, &hash)
	hashUUID(e.WaitingTxnID, &hash)
	hash.Add(uint64(e.CollectionTs.UnixMilli()))
	return hash.Sum()
}

// hashUUID adds the hash of the uuid into the fnv.
// An uuid is a 16 byte array. To hash UUID, we treat it as two uint64 integers,
// since uint64 is 8-byte. This is why we decode the byte array twice and add
// the resulting uint64 into the fnv each time.
func hashUUID(u uuid.UUID, fnv *util.FNV64) {
	b := u.GetBytes()

	b, val, err := encoding.DecodeUint64Descending(b)
	if err != nil {
		panic(err)
	}
	fnv.Add(val)
	_, val, err = encoding.DecodeUint64Descending(b)
	if err != nil {
		panic(err)
	}
	fnv.Add(val)
}
