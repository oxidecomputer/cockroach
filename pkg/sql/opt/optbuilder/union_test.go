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

package optbuilder

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

func TestUnionType(t *testing.T) {
	testCases := []struct {
		left, right, expected *types.T
	}{
		{
			left:     types.Unknown,
			right:    types.Int,
			expected: types.Int,
		},
		{
			left:     types.Int,
			right:    types.Unknown,
			expected: types.Int,
		},
		{
			left:     types.Int4,
			right:    types.Int,
			expected: types.Int,
		},
		{
			left:     types.Int4,
			right:    types.Int2,
			expected: types.Int4,
		},
		{
			left:     types.Float4,
			right:    types.Float,
			expected: types.Float,
		},
		{
			left:     types.MakeDecimal(12 /* precision */, 5 /* scale */),
			right:    types.MakeDecimal(10 /* precision */, 7 /* scale */),
			expected: types.MakeDecimal(10 /* precision */, 7 /* scale */),
		},
		{
			// At the same scale, we use the left type.
			left:     types.MakeDecimal(10 /* precision */, 1 /* scale */),
			right:    types.MakeDecimal(12 /* precision */, 1 /* scale */),
			expected: types.MakeDecimal(10 /* precision */, 1 /* scale */),
		},
		{
			left:     types.Int4,
			right:    types.Decimal,
			expected: types.Decimal,
		},
		{
			left:     types.Decimal,
			right:    types.Float,
			expected: types.Decimal,
		},
		{
			// Error.
			left:     types.Float,
			right:    types.String,
			expected: nil,
		},
		{
			// Error.
			left:     types.MakeArray(types.MakeTuple([]*types.T{types.Any})),
			right:    types.MakeArray(types.MakeTuple([]*types.T{types.Bool})),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		result := func() *types.T {
			defer func() {
				// Swallow any error and return nil.
				_ = recover()
			}()
			return determineUnionType(tc.left, tc.right, "test")
		}()
		toStr := func(t *types.T) string {
			if t == nil {
				return "<nil>"
			}
			return t.SQLString()
		}
		if toStr(result) != toStr(tc.expected) {
			t.Errorf(
				"left: %s  right: %s  expected: %s  got: %s",
				toStr(tc.left), toStr(tc.right), toStr(tc.expected), toStr(result),
			)
		}
	}
}
