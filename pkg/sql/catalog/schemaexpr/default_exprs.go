// Copyright 2016 The Cockroach Authors.
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

	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/transform"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// MakeDefaultExprs returns a slice of the default expressions for the slice
// of input column descriptors, or nil if none of the input column descriptors
// have default expressions.
// The length of the result slice matches the length of the input column descriptors.
// For every column that has no default expression, a NULL expression is reported
// as default.
func MakeDefaultExprs(
	ctx context.Context,
	cols []catalog.Column,
	txCtx *transform.ExprTransformContext,
	evalCtx *tree.EvalContext,
	semaCtx *tree.SemaContext,
) ([]tree.TypedExpr, error) {
	// Check to see if any of the columns have DEFAULT expressions. If there
	// are no DEFAULT expressions, we don't bother with constructing the
	// defaults map as the defaults are all NULL.
	haveDefaults := false
	for _, col := range cols {
		if col.HasDefault() {
			haveDefaults = true
			break
		}
	}
	if !haveDefaults {
		return nil, nil
	}

	// Build the default expressions map from the parsed SELECT statement.
	defaultExprs := make([]tree.TypedExpr, 0, len(cols))
	exprStrings := make([]string, 0, len(cols))
	for _, col := range cols {
		if col.HasDefault() {
			exprStrings = append(exprStrings, col.GetDefaultExpr())
		}
	}
	exprs, err := parser.ParseExprs(exprStrings)
	if err != nil {
		return nil, err
	}

	defExprIdx := 0
	for _, col := range cols {
		if !col.HasDefault() {
			defaultExprs = append(defaultExprs, tree.DNull)
			continue
		}
		expr := exprs[defExprIdx]
		typedExpr, err := tree.TypeCheck(ctx, expr, semaCtx, col.GetType())
		if err != nil {
			return nil, err
		}
		if typedExpr, err = txCtx.NormalizeExpr(evalCtx, typedExpr); err != nil {
			return nil, err
		}
		defaultExprs = append(defaultExprs, typedExpr)
		defExprIdx++
	}
	return defaultExprs, nil
}

// ProcessColumnSet returns columns in cols, and other writable
// columns in tableDesc that fulfills a given criteria in inSet.
func ProcessColumnSet(
	cols []catalog.Column, tableDesc catalog.TableDescriptor, inSet func(column catalog.Column) bool,
) []catalog.Column {
	var colIDSet catalog.TableColSet
	for i := range cols {
		colIDSet.Add(cols[i].GetID())
	}

	// Add all public or columns in DELETE_AND_WRITE_ONLY state
	// that satisfy the condition.
	ret := make([]catalog.Column, 0, len(tableDesc.AllColumns()))
	ret = append(ret, cols...)
	for _, col := range tableDesc.WritableColumns() {
		if inSet(col) {
			if !colIDSet.Contains(col.GetID()) {
				colIDSet.Add(col.GetID())
				ret = append(ret, col)
			}
		}
	}
	return ret
}
