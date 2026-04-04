// Copyright 2017 The Cockroach Authors.
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

package sql

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// analyzeExpr performs semantic analysis of an expression, including:
// - replacing sub-queries by a sql.subquery node;
// - resolving names (optional);
// - type checking (with optional type enforcement);
// - normalization.
// The parameters sources and IndexedVars, if both are non-nil, indicate
// name resolution should be performed. The IndexedVars map will be filled
// as a result.
func (p *planner) analyzeExpr(
	ctx context.Context,
	raw tree.Expr,
	source *colinfo.DataSourceInfo,
	iVarHelper tree.IndexedVarHelper,
	expectedType *types.T,
	requireType bool,
	typingContext string,
) (tree.TypedExpr, error) {
	// Perform optional name resolution.
	resolved := raw
	if source != nil {
		var err error
		resolved, err = p.resolveNames(raw, source, iVarHelper)
		if err != nil {
			return nil, err
		}
	}

	// Type check.
	var typedExpr tree.TypedExpr
	var err error
	p.semaCtx.IVarContainer = iVarHelper.Container()
	if requireType {
		typedExpr, err = tree.TypeCheckAndRequire(ctx, resolved, &p.semaCtx,
			expectedType, typingContext)
	} else {
		typedExpr, err = tree.TypeCheck(ctx, resolved, &p.semaCtx, expectedType)
	}
	p.semaCtx.IVarContainer = nil
	if err != nil {
		return nil, err
	}

	// Normalize.
	return p.txCtx.NormalizeExpr(p.EvalContext(), typedExpr)
}
