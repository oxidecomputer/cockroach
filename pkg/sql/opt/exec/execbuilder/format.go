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

package execbuilder

import (
	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

func init() {
	// Install the interceptor that implements the ExprFmtHideScalars functionality.
	memo.ScalarFmtInterceptor = fmtInterceptor
}

// fmtInterceptor is a function suitable for memo.ScalarFmtInterceptor. It detects
// if an expression tree contains only scalar expressions; if so, it tries to
// execbuild them and print the SQL expressions.
func fmtInterceptor(f *memo.ExprFmtCtx, scalar opt.ScalarExpr) string {
	// An AssignmentCastExpr is built as a crdb_internal.assignment_cast
	// function call by execbuilder. Formatting it as such would be confusing in
	// an opt tree, because it would look like a FunctionExpr. So we print the
	// full nodes instead.
	if !onlyScalarsWithoutAssignmentCasts(scalar) {
		return ""
	}

	switch scalar.Op() {
	case opt.FiltersOp:
		// Let the filters node show up; we will apply the code on each filter.
		return ""
	}

	// Build the scalar expression and format it as a single string.
	bld := New(
		nil, /* factory */
		nil, /* optimizer */
		f.Memo,
		nil, /* catalog */
		scalar,
		nil,   /* evalCtx */
		false, /* allowAutoCommit */
	)
	expr, err := bld.BuildScalar()
	if err != nil {
		// Not all scalar operators are supported (e.g. Projections).
		return ""
	}
	fmtCtx := tree.NewFmtCtx(
		tree.FmtSimple,
		tree.FmtIndexedVarFormat(func(ctx *tree.FmtCtx, idx int) {
			ctx.WriteString(f.ColumnString(opt.ColumnID(idx + 1)))
		}),
	)
	expr.Format(fmtCtx)
	return fmtCtx.String()
}

func onlyScalarsWithoutAssignmentCasts(expr opt.Expr) bool {
	if !opt.IsScalarOp(expr) || expr.Op() == opt.AssignmentCastOp {
		return false
	}
	for i, n := 0, expr.ChildCount(); i < n; i++ {
		if !onlyScalarsWithoutAssignmentCasts(expr.Child(i)) {
			return false
		}
	}
	return true
}
