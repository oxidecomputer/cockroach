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

package memsize

import (
	"time"
	"unsafe"

	"github.com/cockroachdb/apd/v3"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/duration"
)

// These constants are provided to help estimate memory usage to report to
// BoundAccount and BytesMonitor.
const (
	// Bool is the in-memory size of a bool in bytes.
	Bool = int64(unsafe.Sizeof(true))

	// Int is the in-memory size of an int in bytes.
	Int = int64(unsafe.Sizeof(int(0)))

	// Int16 is the in-memory size of an int16 in bytes.
	Int16 = int64(unsafe.Sizeof(int16(0)))

	// Int32 is the in-memory size of an int32 in bytes.
	Int32 = int64(unsafe.Sizeof(int32(0)))

	// Uint32 is the in-memory size of a uint32 in bytes.
	Uint32 = int64(unsafe.Sizeof(uint32(0)))

	// Int64 is the in-memory size of an int64 in bytes.
	Int64 = int64(unsafe.Sizeof(int64(0)))

	// Uint64 is the in-memory size of a uint64 in bytes.
	Uint64 = int64(unsafe.Sizeof(uint64(0)))

	// Float64 is the in-memory size of a float64 in bytes.
	Float64 = int64(unsafe.Sizeof(float64(0)))

	// Time is the in-memory size of a time.Time in bytes.
	Time = int64(unsafe.Sizeof(time.Time{}))

	// Duration is the in-memory size of a duration.Duration in bytes.
	Duration = int64(unsafe.Sizeof(duration.Duration{}))

	// Decimal is the in-memory size of an apd.Decimal in bytes.
	Decimal = int64(unsafe.Sizeof(apd.Decimal{}))

	// String is the in-memory size of an empty string in bytes.
	String = int64(unsafe.Sizeof(""))

	// BoolSliceOverhead is the in-memory overhead of a []bool in bytes.
	BoolSliceOverhead = int64(unsafe.Sizeof([]bool{}))

	// IntSliceOverhead is the in-memory overhead of an []int in bytes.
	IntSliceOverhead = int64(unsafe.Sizeof([]int{}))

	// DatumOverhead is the in-memory overhead of a tree.Datum in bytes.
	DatumOverhead = int64(unsafe.Sizeof(tree.Datum(nil)))

	// DatumsOverhead is the in-memory overhead of a []tree.Datum in bytes.
	DatumsOverhead = int64(unsafe.Sizeof([]tree.Datum{}))

	// RowsOverhead is the in-memory overhead of a [][]tree.Datum in bytes.
	RowsOverhead = int64(unsafe.Sizeof([][]tree.Datum{}))

	// RowsSliceOverhead is the in-memory overhead of a [][][]tree.Datum in
	// bytes.
	RowsSliceOverhead = int64(unsafe.Sizeof([][][]tree.Datum{}))

	// MapEntryOverhead is an estimate of the size of each item in a map in
	// addition to the space occupied by the key and value. This value was
	// determined empirically using runtime.GC() and runtime.ReadMemStats() to
	// analyze the memory used by a map. This overhead appears to be independent
	// of the key and value data types.
	MapEntryOverhead = 64
)
