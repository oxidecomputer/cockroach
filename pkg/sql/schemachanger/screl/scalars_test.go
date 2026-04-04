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

package screl

import (
	"math"
	"reflect"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/catid"
	types "github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/stretchr/testify/require"
)

// TestAllElementsHaveDescID ensures that all element types have a DescID.
func TestAllElementsHaveDescID(t *testing.T) {
	typ := reflect.TypeOf((*scpb.ElementProto)(nil)).Elem()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		elem := reflect.New(f.Type.Elem()).Interface().(scpb.Element)
		require.Equal(t, descpb.ID(0), GetDescID(elem))
	}
}

func TestAllDescIDsAndContainsDescID(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    scpb.Element
		expected []catid.DescID
	}{
		{
			name: "schema parent",
			input: &scpb.SchemaParent{
				SchemaID:         1,
				ParentDatabaseID: 2,
			},
			expected: []catid.DescID{1, 2},
		},
		{
			name: "default expr",
			input: &scpb.ColumnDefaultExpression{
				TableID:  1,
				ColumnID: 10,
				Expression: scpb.Expression{
					Expr:            "foo",
					UsesTypeIDs:     []catid.DescID{2, 3},
					UsesSequenceIDs: []catid.DescID{4, 5},
				},
			},
			expected: []catid.DescID{1, 2, 3, 4, 5},
		},
		{
			name: "udf column",
			input: &scpb.ColumnType{
				TableID:  1,
				ColumnID: 10,
				TypeT: scpb.TypeT{
					Type:          types.Any,
					ClosedTypeIDs: []catid.DescID{2, 3},
				},
			},
			expected: []catid.DescID{1, 2, 3},
		},
		{
			name: "computed column",
			input: &scpb.ColumnType{
				TableID:  1,
				ColumnID: 10,
				TypeT: scpb.TypeT{
					Type:          types.Any,
					ClosedTypeIDs: []catid.DescID{2, 3},
				},
				ComputeExpr: &scpb.Expression{
					Expr:            "foo",
					UsesTypeIDs:     []catid.DescID{3, 4},
					UsesSequenceIDs: []catid.DescID{5, 6},
				},
			},
			expected: []catid.DescID{1, 2, 3, 4, 5, 6},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.ElementsMatch(t, tc.expected, AllDescIDs(tc.input).Ordered())
			for _, id := range tc.expected {
				require.Truef(t, ContainsDescID(tc.input, id), "contains %d", id)
			}
			require.False(t, ContainsDescID(tc.input, 0))
			require.False(t, ContainsDescID(tc.input, math.MaxUint32))
		})
	}
}
