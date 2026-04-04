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

package testutils

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/cat"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/norm"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/optbuilder"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/xform"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// BuildQuery initializes an optimizer and builds the given sql statement.
func BuildQuery(
	t *testing.T, o *xform.Optimizer, catalog cat.Catalog, evalCtx *tree.EvalContext, sql string,
) {
	stmt, err := parser.ParseOne(sql)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	semaCtx := tree.MakeSemaContext()
	if err := semaCtx.Placeholders.Init(stmt.NumPlaceholders, nil /* typeHints */); err != nil {
		t.Fatal(err)
	}
	semaCtx.Annotations = tree.MakeAnnotations(stmt.NumAnnotations)
	o.Init(evalCtx, catalog)
	err = optbuilder.New(ctx, &semaCtx, evalCtx, catalog, o.Factory(), stmt.AST).Build()
	if err != nil {
		t.Fatal(err)
	}
}

// BuildScalar builds the given input string as a ScalarExpr and returns it.
func BuildScalar(
	t *testing.T, f *norm.Factory, semaCtx *tree.SemaContext, evalCtx *tree.EvalContext, input string,
) opt.ScalarExpr {
	expr, err := parser.ParseExpr(input)
	if err != nil {
		t.Fatal(err)
	}

	b := optbuilder.NewScalar(context.Background(), semaCtx, evalCtx, f)
	if err := b.Build(expr); err != nil {
		t.Fatal(err)
	}

	return f.Memo().RootExpr().(opt.ScalarExpr)
}

// BuildFilters builds the given input string as a FiltersExpr and returns it.
// Calls a subset of the normalization rules that would apply if these filters
// were built as part of a Select or Join.
func BuildFilters(
	t *testing.T, f *norm.Factory, semaCtx *tree.SemaContext, evalCtx *tree.EvalContext, input string,
) memo.FiltersExpr {
	if input == "" {
		return memo.TrueFilter
	}
	root := BuildScalar(t, f, semaCtx, evalCtx, input)

	if _, ok := root.(*memo.TrueExpr); ok {
		return memo.TrueFilter
	}
	filters := memo.FiltersExpr{f.ConstructFiltersItem(root)}
	filters = f.CustomFuncs().SimplifyFilters(filters)
	filters = f.CustomFuncs().ConsolidateFilters(filters)
	return filters
}
