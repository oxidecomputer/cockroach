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

package workload_test

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/testutils/skip"
	"github.com/cockroachdb/cockroach/pkg/util/bufalloc"
	"github.com/cockroachdb/cockroach/pkg/workload"
	"github.com/cockroachdb/cockroach/pkg/workload/bank"
	"github.com/cockroachdb/cockroach/pkg/workload/tpcc"
	"github.com/cockroachdb/cockroach/pkg/workload/tpch"
)

func columnByteSize(col coldata.Vec) int64 {
	switch t := col.Type(); col.CanonicalTypeFamily() {
	case types.IntFamily:
		switch t.Width() {
		case 0, 64:
			return int64(len(col.Int64()) * 8)
		case 16:
			return int64(len(col.Int16()) * 2)
		default:
			panic(fmt.Sprintf("unexpected int width: %d", t.Width()))
		}
	case types.FloatFamily:
		return int64(len(col.Float64()) * 8)
	case types.BytesFamily:
		// We subtract the overhead to be in line with Int64 and Float64 cases.
		return col.Bytes().Size() - coldata.FlatBytesOverhead
	default:
		panic(fmt.Sprintf(`unhandled type %s`, t))
	}
}

func benchmarkInitialData(b *testing.B, gen workload.Generator) {
	tables := gen.Tables()

	var bytes int64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Share the Batch and ByteAllocator across tables but not across benchmark
		// iterations.
		cb := coldata.NewMemBatch(nil /* types */, coldata.StandardColumnFactory)
		var a bufalloc.ByteAllocator
		for _, table := range tables {
			for rowIdx := 0; rowIdx < table.InitialRows.NumBatches; rowIdx++ {
				a = a.Truncate()
				table.InitialRows.FillBatch(rowIdx, cb, &a)
				for _, col := range cb.ColVecs() {
					bytes += columnByteSize(col)
				}
			}
		}
	}
	b.StopTimer()
	b.SetBytes(bytes / int64(b.N))
}

func BenchmarkInitialData(b *testing.B) {
	b.Run(`tpcc/warehouses=1`, func(b *testing.B) {
		benchmarkInitialData(b, tpcc.FromWarehouses(1))
	})
	b.Run(`bank/rows=1000`, func(b *testing.B) {
		benchmarkInitialData(b, bank.FromRows(1000))
	})
	b.Run(`tpch/scaleFactor=1`, func(b *testing.B) {
		skip.UnderShort(b, "tpch loads a lot of data")
		benchmarkInitialData(b, tpch.FromScaleFactor(1))
	})
}
