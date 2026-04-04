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

package colinfo

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// CheckDatumTypeFitsColumnType verifies that a given scalar value
// type is valid to be stored in a column of the given column type.
//
// For the purpose of this analysis, column type aliases are not
// considered to be different (eg. TEXT and VARCHAR will fit the same
// scalar type String).
//
// This is used by the UPDATE, INSERT and UPSERT code.
func CheckDatumTypeFitsColumnType(col catalog.Column, typ *types.T) error {
	if typ.Family() == types.UnknownFamily {
		return nil
	}
	if !typ.Equivalent(col.GetType()) {
		return pgerror.Newf(pgcode.DatatypeMismatch,
			"value type %s doesn't match type %s of column %q",
			typ.String(), col.GetType().String(), tree.ErrNameString(col.GetName()))
	}
	return nil
}

// CanHaveCompositeKeyEncoding returns true if key columns of the given kind can
// have a composite encoding. For such types, it can be decided on a
// case-by-base basis whether a given Datum requires the composite encoding.
//
// As an example of a composite encoding, collated string key columns are
// encoded partly as a key and partly as a value. The key part is the collation
// key, so that different strings that collate equal cannot both be used as
// keys. The value part is the usual UTF-8 encoding of the string, stored so
// that it can be recovered later for inspection/display.
func CanHaveCompositeKeyEncoding(typ *types.T) bool {
	switch typ.Family() {
	case types.FloatFamily,
		types.DecimalFamily,
		types.CollatedStringFamily:
		return true
	case types.ArrayFamily:
		return CanHaveCompositeKeyEncoding(typ.ArrayContents())
	case types.TupleFamily:
		for _, t := range typ.TupleContents() {
			if CanHaveCompositeKeyEncoding(t) {
				return true
			}
		}
		return false
	case types.BoolFamily,
		types.IntFamily,
		types.DateFamily,
		types.TimestampFamily,
		types.IntervalFamily,
		types.StringFamily,
		types.BytesFamily,
		types.TimestampTZFamily,
		types.OidFamily,
		types.UuidFamily,
		types.INetFamily,
		types.TimeFamily,
		types.JsonFamily,
		types.TimeTZFamily,
		types.BitFamily,
		types.GeometryFamily,
		types.GeographyFamily,
		types.EnumFamily,
		types.Box2DFamily:
		return false
	case types.UnknownFamily,
		types.AnyFamily:
		fallthrough
	default:
		return true
	}
}
