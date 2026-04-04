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

package testcat

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/opt/cat"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/errors"
)

// typeCheckTableExpr takes an expression that refers to table columns by name
// and returns its resulting type.
func typeCheckTableExpr(e tree.Expr, cols []cat.Column) *types.T {
	tableCols := make(map[tree.Name]*tableCol)
	for i := range cols {
		if cols[i].Kind() == cat.Ordinary {
			tableCols[cols[i].ColName()] = &tableCol{
				name: cols[i].ColName(),
				typ:  cols[i].DatumType(),
			}
		}
	}
	// To type check the expression, we resolve column name by replacing them with
	// *tableCol objects.
	visitFn := func(expr tree.Expr) (recurse bool, newExpr tree.Expr, err error) {
		if n, ok := expr.(*tree.UnresolvedName); ok {
			v, err := n.NormalizeVarName()
			if err != nil {
				return false, nil, err
			}
			item, ok := v.(*tree.ColumnItem)
			if !ok {
				return false, nil, errors.Newf("unsupported column selector %s", n)
			}
			col, ok := tableCols[item.ColumnName]
			if !ok {
				return false, nil, errors.Newf("unknown column %s", item.ColumnName)
			}
			return false, col, nil
		}
		return true, expr, nil
	}
	resolved, err := tree.SimpleVisit(e, visitFn)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	semaCtx := tree.MakeSemaContext()
	typedExpr, err := resolved.TypeCheck(ctx, &semaCtx, types.Any)
	if err != nil {
		panic(err)
	}
	return typedExpr.ResolvedType()
}

// tableCol is a tree.Expr node that replaces unresolved columns (to allow
// type-checking).
type tableCol struct {
	name tree.Name
	typ  *types.T
}

var _ tree.TypedExpr = &tableCol{}
var _ tree.VariableExpr = &tableCol{}

func (c *tableCol) String() string {
	return tree.AsString(c)
}

// Format implements the NodeFormatter interface.
func (c *tableCol) Format(ctx *tree.FmtCtx) {
	ctx.FormatNode(&c.name)
}

// Walk is part of the tree.Expr interface.
func (c *tableCol) Walk(v tree.Visitor) tree.Expr {
	return c
}

// TypeCheck is part of the tree.Expr interface.
func (c *tableCol) TypeCheck(
	_ context.Context, _ *tree.SemaContext, desired *types.T,
) (tree.TypedExpr, error) {
	return c, nil
}

// ResolvedType is part of the tree.TypedExpr interface.
func (c *tableCol) ResolvedType() *types.T {
	return c.typ
}

// Eval is part of the tree.TypedExpr interface.
func (*tableCol) Eval(_ *tree.EvalContext) (tree.Datum, error) {
	panic("not implemented")
}

// Variable is part of the tree.VariableExpr interface. This prevents the
// column from being evaluated during normalization.
func (*tableCol) Variable() {}
