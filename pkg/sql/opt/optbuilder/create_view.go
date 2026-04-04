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

package optbuilder

import (
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlerrors"
	"github.com/cockroachdb/cockroach/pkg/util"
)

func (b *Builder) buildCreateView(cv *tree.CreateView, inScope *scope) (outScope *scope) {
	b.DisableMemoReuse = true
	sch, resName := b.resolveSchemaForCreate(&cv.Name)
	schID := b.factory.Metadata().AddSchema(sch)
	viewName := tree.MakeTableNameFromPrefix(resName, tree.Name(cv.Name.Object()))

	// We build the select statement to:
	//  - check the statement semantically,
	//  - get the fully resolved names into the AST, and
	//  - collect the view dependencies in b.viewDeps.
	// The result is not otherwise used.
	b.insideViewDef = true
	b.trackViewDeps = true
	b.qualifyDataSourceNamesInAST = true
	if b.sourceViews == nil {
		b.sourceViews = make(map[string]struct{})
	}
	b.sourceViews[viewName.FQString()] = struct{}{}
	defer func() {
		b.insideViewDef = false
		b.trackViewDeps = false
		b.viewDeps = nil
		b.viewTypeDeps = util.FastIntSet{}
		b.qualifyDataSourceNamesInAST = false
		delete(b.sourceViews, viewName.FQString())
	}()

	defScope := b.buildStmtAtRoot(cv.AsSource, nil /* desiredTypes */)

	p := defScope.makePhysicalProps().Presentation
	if len(cv.ColumnNames) != 0 {
		if len(p) != len(cv.ColumnNames) {
			panic(sqlerrors.NewSyntaxErrorf(
				"CREATE VIEW specifies %d column name%s, but data source has %d column%s",
				len(cv.ColumnNames), util.Pluralize(int64(len(cv.ColumnNames))),
				len(p), util.Pluralize(int64(len(p)))),
			)
		}
		// Override the columns.
		for i := range p {
			p[i].Alias = string(cv.ColumnNames[i])
		}
	}

	// If the type of any column that this view references is user
	// defined, add a type dependency between this view and the UDT.
	if b.trackViewDeps {
		for _, d := range b.viewDeps {
			if !d.ColumnOrdinals.Empty() {
				d.ColumnOrdinals.ForEach(func(ord int) {
					ids, err := d.DataSource.CollectTypes(ord)
					if err != nil {
						panic(err)
					}
					for _, id := range ids {
						b.viewTypeDeps.Add(int(id))
					}
				})
			}
		}
	}

	outScope = b.allocScope()
	outScope.expr = b.factory.ConstructCreateView(
		&memo.CreateViewPrivate{
			Schema:       schID,
			ViewName:     &viewName,
			IfNotExists:  cv.IfNotExists,
			Replace:      cv.Replace,
			Persistence:  cv.Persistence,
			Materialized: cv.Materialized,
			ViewQuery:    tree.AsStringWithFlags(cv.AsSource, tree.FmtParsable),
			Columns:      p,
			Deps:         b.viewDeps,
			TypeDeps:     b.viewTypeDeps,
		},
	)
	return outScope
}
