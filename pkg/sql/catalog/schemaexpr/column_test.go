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

package schemaexpr

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

func TestDequalifyColumnRefs(t *testing.T) {
	ctx := context.Background()

	database := tree.Name("foo")
	table := tree.Name("bar")
	tn := tree.MakeTableNameWithSchema(database, tree.PublicSchemaName, table)

	cols := []descpb.ColumnDescriptor{
		{Name: "a", Type: types.Int},
		{Name: "b", Type: types.String},
	}

	testData := []struct {
		expr     string
		expected string
	}{
		{"a", "a"},
		{"bar.a", "a"},
		{"foo.bar.a", "a"},
		{"a > 0", "a > 0"},
		{"bar.a > 0", "a > 0"},
		{"foo.bar.a > 0", "a > 0"},
		{"a > 0 AND b = 'baz'", "(a > 0) AND (b = 'baz')"},
		{"bar.a > 0 AND bar.b = 'baz'", "(a > 0) AND (b = 'baz')"},
		{"foo.bar.a > 0 AND foo.bar.b = 'baz'", "(a > 0) AND (b = 'baz')"},
		{"bar.a > 0 AND b = 'baz'", "(a > 0) AND (b = 'baz')"},
		{"foo.bar.a > 0 AND b = 'baz'", "(a > 0) AND (b = 'baz')"},
		{"a > 0 AND bar.b = 'baz'", "(a > 0) AND (b = 'baz')"},
		{"a > 0 AND foo.bar.b = 'baz'", "(a > 0) AND (b = 'baz')"},
	}

	for _, d := range testData {
		t.Run(d.expr, func(t *testing.T) {
			expr, err := parser.ParseExpr(d.expr)
			if err != nil {
				t.Fatalf("%s: unexpected error: %s", d.expr, err)
			}

			source := colinfo.NewSourceInfoForSingleTable(
				tn, colinfo.ResultColumnsFromColDescs(
					descpb.ID(1),
					len(cols),
					func(i int) *descpb.ColumnDescriptor {
						return &cols[i]
					},
				),
			)

			deqExpr, err := DequalifyColumnRefs(ctx, source, expr)
			if err != nil {
				t.Fatalf("%s: expected success, but found error: %s", d.expr, err)
			}

			if deqExpr != d.expected {
				t.Errorf("%s: expected %q, got %q", d.expr, d.expected, deqExpr)
			}
		})
	}
}

func TestRenameColumn(t *testing.T) {
	from := tree.Name("foo")
	to := tree.Name("bar")

	testData := []struct {
		expr     string
		expected string
	}{
		{"foo", "bar"},
		{"foo = 1", "bar = 1"},
		{"foo = 1 AND baz = 3", "(bar = 1) AND (baz = 3)"},
		{"baz = 3 OR foo = 1", "(baz = 3) OR (bar = 1)"},
		{"timezone(baz, foo::TIMESTAMPTZ) > now()", "timezone(baz, bar::TIMESTAMPTZ) > now()"},
	}

	for _, d := range testData {
		t.Run(d.expr, func(t *testing.T) {
			res, err := RenameColumn(d.expr, from, to)
			if err != nil {
				t.Fatalf("%s: unexpected error: %s", d.expr, err)
			}

			if res != d.expected {
				t.Errorf("%s: expected %q, got %q", d.expr, d.expected, res)
			}
		})
	}
}
