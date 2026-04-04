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

package colexeccmp

import "github.com/cockroachdb/cockroach/pkg/sql/sem/tree"

// ComparisonExprAdapter is a utility interface that is implemented by several
// structs that behave as an adapter from tree.ComparisonExpr to a vectorized
// friendly model.
type ComparisonExprAdapter interface {
	Eval(left, right tree.Datum) (tree.Datum, error)
}

// NewComparisonExprAdapter returns a new ComparisonExprAdapter for the provided
// expression.
func NewComparisonExprAdapter(
	expr *tree.ComparisonExpr, evalCtx *tree.EvalContext,
) ComparisonExprAdapter {
	base := cmpExprAdapterBase{
		fn:      expr.Fn.Fn,
		evalCtx: evalCtx,
	}
	op := expr.Operator
	if op.Symbol.HasSubOperator() {
		return &cmpWithSubOperatorExprAdapter{
			cmpExprAdapterBase: base,
			expr:               expr,
		}
	}
	nullable := expr.Fn.NullableArgs
	_, _, _, flipped, negate := tree.FoldComparisonExpr(op, nil /* left */, nil /* right */)
	if nullable {
		if flipped {
			if negate {
				return &cmpNullableFlippedNegateExprAdapter{cmpExprAdapterBase: base}
			}
			return &cmpNullableFlippedExprAdapter{cmpExprAdapterBase: base}
		}
		if negate {
			return &cmpNullableNegateExprAdapter{cmpExprAdapterBase: base}
		}
		return &cmpNullableExprAdapter{cmpExprAdapterBase: base}
	}
	if flipped {
		if negate {
			return &cmpFlippedNegateExprAdapter{cmpExprAdapterBase: base}
		}
		return &cmpFlippedExprAdapter{cmpExprAdapterBase: base}
	}
	if negate {
		return &cmpNegateExprAdapter{cmpExprAdapterBase: base}
	}
	return &cmpExprAdapter{cmpExprAdapterBase: base}
}
