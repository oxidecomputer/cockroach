// Copyright 2020 The Cockroach Authors.
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

package typeconv

import (
	"fmt"
	"time"

	"github.com/cockroachdb/apd/v3"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/duration"
	"github.com/cockroachdb/cockroach/pkg/util/json"
)

// DatumVecCanonicalTypeFamily is the "canonical" type family of all types that
// are physically represented by coldata.DatumVec.
var DatumVecCanonicalTypeFamily = types.Family(1000000)

// TypeFamilyToCanonicalTypeFamily converts all type families to their
// "canonical" counterparts. "Canonical" type families are representatives
// from a set of "equivalent" type families where "equivalence" means having
// the same physical representation.
//
// All type families that do not have an optimized physical representation are
// handled by using tree.Datums, and such types are mapped to
// DatumVecCanonicalTypeFamily.
func TypeFamilyToCanonicalTypeFamily(family types.Family) types.Family {
	switch family {
	case types.BoolFamily:
		return types.BoolFamily
	case types.BytesFamily, types.StringFamily, types.UuidFamily:
		return types.BytesFamily
	case types.DecimalFamily:
		return types.DecimalFamily
	case types.JsonFamily:
		return types.JsonFamily
	case types.IntFamily, types.DateFamily:
		return types.IntFamily
	case types.FloatFamily:
		return types.FloatFamily
	case types.TimestampTZFamily, types.TimestampFamily:
		return types.TimestampTZFamily
	case types.IntervalFamily:
		return types.IntervalFamily
	default:
		// TODO(yuzefovich): consider adding native support for
		// types.UnknownFamily.
		return DatumVecCanonicalTypeFamily
	}
}

// ToCanonicalTypeFamilies converts typs to the corresponding canonical type
// families.
func ToCanonicalTypeFamilies(typs []*types.T) []types.Family {
	families := make([]types.Family, len(typs))
	for i := range typs {
		families[i] = TypeFamilyToCanonicalTypeFamily(typs[i].Family())
	}
	return families
}

// UnsafeFromGoType returns the type for a Go value, if applicable. Shouldn't
// be used at runtime. This method is unsafe because multiple logical types can
// be represented by the same physical type. Types that are backed by DatumVec
// are *not* supported by this function.
func UnsafeFromGoType(v interface{}) *types.T {
	switch t := v.(type) {
	case int16:
		return types.Int2
	case int32:
		return types.Int4
	case int, int64:
		return types.Int
	case bool:
		return types.Bool
	case float64:
		return types.Float
	case []byte:
		return types.Bytes
	case string:
		return types.String
	case apd.Decimal:
		return types.Decimal
	case time.Time:
		return types.TimestampTZ
	case duration.Duration:
		return types.Interval
	case json.JSON:
		return types.Jsonb
	default:
		panic(fmt.Sprintf("type %s not supported yet", t))
	}
}

// TypesSupportedNatively contains types that are supported natively by the
// vectorized engine.
var TypesSupportedNatively []*types.T

func init() {
	for _, t := range types.Scalar {
		if TypeFamilyToCanonicalTypeFamily(t.Family()) == DatumVecCanonicalTypeFamily {
			continue
		}
		if t.Family() == types.IntFamily {
			TypesSupportedNatively = append(TypesSupportedNatively, types.Int2)
			TypesSupportedNatively = append(TypesSupportedNatively, types.Int4)
		}
		TypesSupportedNatively = append(TypesSupportedNatively, t)
	}
}
