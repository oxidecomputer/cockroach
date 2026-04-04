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

package colexecargs

import (
	"sync"

	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

var exprHelperPool = sync.Pool{
	New: func() interface{} {
		return &ExprHelper{}
	},
}

// NewExprHelper returns a new ExprHelper.
func NewExprHelper() *ExprHelper {
	return exprHelperPool.Get().(*ExprHelper)
}

// ExprHelper is a utility struct that helps with expression handling in the
// vectorized engine.
type ExprHelper struct {
	helper  execinfrapb.ExprHelper
	SemaCtx *tree.SemaContext
}

// ProcessExpr processes the given expression and returns a well-typed
// expression. Note that SemaCtx must be already set on h.
//
// evalCtx will not be mutated.
func (h *ExprHelper) ProcessExpr(
	expr execinfrapb.Expression, evalCtx *tree.EvalContext, typs []*types.T,
) (tree.TypedExpr, error) {
	if expr.LocalExpr != nil {
		return expr.LocalExpr, nil
	}
	h.helper.Types = typs
	tempVars := tree.MakeIndexedVarHelper(&h.helper, len(typs))
	return execinfrapb.DeserializeExpr(expr.Expr, h.SemaCtx, evalCtx, &tempVars)
}
