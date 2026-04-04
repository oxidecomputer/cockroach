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

package colencoding

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/rowenc/valueside"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/duration"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
)

// DecodeTableValueToCol decodes a value encoded by EncodeTableValue, writing
// the result to the rowIdx'th position of the vecIdx'th vector in
// coldata.TypedVecs.
// See the analog in rowenc/column_type_encoding.go.
func DecodeTableValueToCol(
	da *tree.DatumAlloc,
	vecs *coldata.TypedVecs,
	vecIdx int,
	rowIdx int,
	typ encoding.Type,
	dataOffset int,
	valTyp *types.T,
	buf []byte,
) ([]byte, error) {
	// NULL is special because it is a valid value for any type.
	if typ == encoding.Null {
		vecs.Nulls[vecIdx].SetNull(rowIdx)
		return buf[dataOffset:], nil
	}

	// Bool is special because the value is stored in the value tag, so we have
	// to keep the reference to the original slice.
	origBuf := buf
	buf = buf[dataOffset:]

	// Find the position of the target vector among the typed columns of its
	// type.
	colIdx := vecs.ColsMap[vecIdx]

	var err error
	switch valTyp.Family() {
	case types.BoolFamily:
		var b bool
		// A boolean's value is encoded in its tag directly, so we don't have an
		// "Untagged" version of this function. Note that we also use the
		// original buffer.
		buf, b, err = encoding.DecodeBoolValue(origBuf)
		vecs.BoolCols[colIdx][rowIdx] = b
	case types.BytesFamily, types.StringFamily:
		var data []byte
		buf, data, err = encoding.DecodeUntaggedBytesValue(buf)
		vecs.BytesCols[colIdx].Set(rowIdx, data)
	case types.DateFamily:
		var i int64
		buf, i, err = encoding.DecodeUntaggedIntValue(buf)
		vecs.Int64Cols[colIdx][rowIdx] = i
	case types.DecimalFamily:
		buf, err = encoding.DecodeIntoUntaggedDecimalValue(&vecs.DecimalCols[colIdx][rowIdx], buf)
	case types.FloatFamily:
		var f float64
		buf, f, err = encoding.DecodeUntaggedFloatValue(buf)
		vecs.Float64Cols[colIdx][rowIdx] = f
	case types.IntFamily:
		var i int64
		buf, i, err = encoding.DecodeUntaggedIntValue(buf)
		switch valTyp.Width() {
		case 16:
			vecs.Int16Cols[colIdx][rowIdx] = int16(i)
		case 32:
			vecs.Int32Cols[colIdx][rowIdx] = int32(i)
		default:
			// Pre-2.1 BIT was using INT encoding with arbitrary sizes.
			// We map these to 64-bit INT now. See #34161.
			vecs.Int64Cols[colIdx][rowIdx] = i
		}
	case types.JsonFamily:
		var data []byte
		buf, data, err = encoding.DecodeUntaggedBytesValue(buf)
		vecs.JSONCols[colIdx].Bytes.Set(rowIdx, data)
	case types.UuidFamily:
		var data uuid.UUID
		buf, data, err = encoding.DecodeUntaggedUUIDValue(buf)
		// TODO(yuzefovich): we could peek inside the encoding package to skip a
		// couple of conversions.
		if err != nil {
			return buf, err
		}
		vecs.BytesCols[colIdx].Set(rowIdx, data.GetBytes())
	case types.TimestampFamily, types.TimestampTZFamily:
		var t time.Time
		buf, t, err = encoding.DecodeUntaggedTimeValue(buf)
		vecs.TimestampCols[colIdx][rowIdx] = t
	case types.IntervalFamily:
		var d duration.Duration
		buf, d, err = encoding.DecodeUntaggedDurationValue(buf)
		vecs.IntervalCols[colIdx][rowIdx] = d
	// Types backed by tree.Datums.
	default:
		var d tree.Datum
		d, buf, err = valueside.DecodeUntaggedDatum(da, valTyp, buf)
		if err != nil {
			return buf, err
		}
		vecs.DatumCols[colIdx].Set(rowIdx, d)
	}
	return buf, err
}
