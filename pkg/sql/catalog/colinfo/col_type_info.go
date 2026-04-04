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

package colinfo

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/errorutil/unimplemented"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq/oid"
	"golang.org/x/text/language"
)

// ColTypeInfo is a type that allows multiple representations of column type
// information (to avoid conversions and allocations).
type ColTypeInfo struct {
	// Only one of these fields can be set.
	resCols  ResultColumns
	colTypes []*types.T
}

// ColTypeInfoFromResCols creates a ColTypeInfo from ResultColumns.
func ColTypeInfoFromResCols(resCols ResultColumns) ColTypeInfo {
	return ColTypeInfo{resCols: resCols}
}

// ColTypeInfoFromColTypes creates a ColTypeInfo from []ColumnType.
func ColTypeInfoFromColTypes(colTypes []*types.T) ColTypeInfo {
	return ColTypeInfo{colTypes: colTypes}
}

// ColTypeInfoFromColumns creates a ColTypeInfo from []catalog.Column.
func ColTypeInfoFromColumns(columns []catalog.Column) ColTypeInfo {
	colTypes := make([]*types.T, len(columns))
	for i, col := range columns {
		colTypes[i] = col.GetType()
	}
	return ColTypeInfoFromColTypes(colTypes)
}

// NumColumns returns the number of columns in the type.
func (ti ColTypeInfo) NumColumns() int {
	if ti.resCols != nil {
		return len(ti.resCols)
	}
	return len(ti.colTypes)
}

// Type returns the datum type of the i-th column.
func (ti ColTypeInfo) Type(idx int) *types.T {
	if ti.resCols != nil {
		return ti.resCols[idx].Typ
	}
	return ti.colTypes[idx]
}

// ValidateColumnDefType returns an error if the type of a column definition is
// not valid. It is checked when a column is created or altered.
func ValidateColumnDefType(t *types.T) error {
	switch t.Family() {
	case types.StringFamily, types.CollatedStringFamily:
		if t.Family() == types.CollatedStringFamily {
			if _, err := language.Parse(t.Locale()); err != nil {
				return pgerror.Newf(pgcode.Syntax, `invalid locale %s`, t.Locale())
			}
		}

	case types.DecimalFamily:
		switch {
		case t.Precision() == 0 && t.Scale() > 0:
			// TODO (seif): Find right range for error message.
			return errors.New("invalid NUMERIC precision 0")
		case t.Precision() < t.Scale():
			return fmt.Errorf("NUMERIC scale %d must be between 0 and precision %d",
				t.Scale(), t.Precision())
		}

	case types.ArrayFamily:
		if t.ArrayContents().Family() == types.ArrayFamily {
			// Nested arrays are not supported as a column type.
			return errors.Errorf("nested array unsupported as column type: %s", t.String())
		}
		if t.ArrayContents().Family() == types.JsonFamily {
			// JSON arrays are not supported as a column type.
			return unimplemented.NewWithIssueDetailf(23468, t.String(),
				"arrays of JSON unsupported as column type")
		}
		if err := types.CheckArrayElementType(t.ArrayContents()); err != nil {
			return err
		}
		return ValidateColumnDefType(t.ArrayContents())

	case types.BitFamily, types.IntFamily, types.FloatFamily, types.BoolFamily, types.BytesFamily, types.DateFamily,
		types.INetFamily, types.IntervalFamily, types.JsonFamily, types.OidFamily, types.TimeFamily,
		types.TimestampFamily, types.TimestampTZFamily, types.UuidFamily, types.TimeTZFamily,
		types.GeographyFamily, types.GeometryFamily, types.EnumFamily, types.Box2DFamily:
		// These types are OK.

	default:
		return pgerror.Newf(pgcode.InvalidTableDefinition,
			"value type %s cannot be used for table columns", t.String())
	}

	return nil
}

// ColumnTypeIsIndexable returns whether the type t is valid as an indexed column.
func ColumnTypeIsIndexable(t *types.T) bool {
	if t.IsAmbiguous() || t.Family() == types.TupleFamily {
		return false
	}
	// Some inverted index types also have a key encoding, but we don't
	// want to support those yet. See #50659.
	return !MustBeValueEncoded(t) && !ColumnTypeIsInvertedIndexable(t)
}

// ColumnTypeIsInvertedIndexable returns whether the type t is valid to be indexed
// using an inverted index.
func ColumnTypeIsInvertedIndexable(t *types.T) bool {
	if t.IsAmbiguous() || t.Family() == types.TupleFamily {
		return false
	}
	family := t.Family()
	return family == types.JsonFamily || family == types.ArrayFamily ||
		family == types.GeographyFamily || family == types.GeometryFamily
}

// MustBeValueEncoded returns true if columns of the given kind can only be value
// encoded.
func MustBeValueEncoded(semanticType *types.T) bool {
	switch semanticType.Family() {
	case types.ArrayFamily:
		switch semanticType.Oid() {
		case oid.T_int2vector, oid.T_oidvector:
			return true
		default:
			return MustBeValueEncoded(semanticType.ArrayContents())
		}
	case types.JsonFamily, types.TupleFamily, types.GeographyFamily, types.GeometryFamily:
		return true
	}
	return false
}
