// Copyright 2018 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/props/physical"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// buildAlterTableRelocate builds an ALTER RANGE RELOCATE (LEASE).
func (b *Builder) buildAlterRangeRelocate(
	relocate *tree.RelocateRange, inScope *scope,
) (outScope *scope) {

	if err := b.catalog.RequireAdminRole(b.ctx, "ALTER RANGE RELOCATE"); err != nil {
		panic(err)
	}

	// Disable optimizer caching, as we do for other ALTER statements.
	b.DisableMemoReuse = true

	outScope = inScope.push()
	b.synthesizeResultColumns(outScope, colinfo.AlterTableRelocateColumns)

	cmdName := "RELOCATE " + relocate.SubjectReplicas.String()
	colNames := []string{"range ids"}
	colTypes := []*types.T{types.Int}

	outScope = inScope.push()
	b.synthesizeResultColumns(outScope, colinfo.AlterRangeRelocateColumns)

	// We don't allow the input statement to reference outer columns, so we
	// pass a "blank" scope rather than inScope.
	emptyScope := b.allocScope()
	inputScope := b.buildStmt(relocate.Rows, colTypes, emptyScope)
	checkInputColumns(cmdName, inputScope, colNames, colTypes, 1)

	var toStoreID, fromStoreID opt.ScalarExpr
	{
		emptyScope.context = exprKindStoreID
		// We need to save and restore the previous value of the field in
		// semaCtx in case we are recursively called within a subquery
		// context.
		defer b.semaCtx.Properties.Restore(b.semaCtx.Properties)
		b.semaCtx.Properties.Require(emptyScope.context.String(), tree.RejectSpecial)

		toStoreIDExpr := emptyScope.resolveType(relocate.ToStoreID, types.Int)
		toStoreID = b.buildScalar(toStoreIDExpr, emptyScope, nil, nil, nil)
		fromStoreIDExpr := emptyScope.resolveType(relocate.FromStoreID, types.Int)
		fromStoreID = b.buildScalar(fromStoreIDExpr, emptyScope, nil, nil, nil)
	}

	outScope.expr = b.factory.ConstructAlterRangeRelocate(
		inputScope.expr,
		toStoreID,
		fromStoreID,
		&memo.AlterRangeRelocatePrivate{
			SubjectReplicas: relocate.SubjectReplicas,
			Columns:         colsToColList(outScope.cols),
			Props:           physical.MinRequired,
		},
	)
	return outScope
}
