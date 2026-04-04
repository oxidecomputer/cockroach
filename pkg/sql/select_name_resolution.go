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
//
// This file implements the select code that deals with column references
// and resolving column names in expressions.

package sql

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/schemaexpr"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// resolveNames walks the provided expression and resolves all names
// using the tableInfo and iVarHelper.
func (p *planner) resolveNames(
	expr tree.Expr, source *colinfo.DataSourceInfo, ivarHelper tree.IndexedVarHelper,
) (tree.Expr, error) {
	if expr == nil {
		return nil, nil
	}
	return schemaexpr.ResolveNamesUsingVisitor(&p.nameResolutionVisitor, expr, source, ivarHelper, p.SessionData().SearchPath)
}
