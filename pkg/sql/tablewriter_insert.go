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

package sql

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/row"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// tableInserter handles writing kvs and forming table rows for inserts.
type tableInserter struct {
	tableWriterBase
	ri row.Inserter
}

var _ tableWriter = &tableInserter{}

// desc is part of the tableWriter interface.
func (*tableInserter) desc() string { return "inserter" }

// init is part of the tableWriter interface.
func (ti *tableInserter) init(
	_ context.Context, txn *kv.Txn, evalCtx *tree.EvalContext, sv *settings.Values,
) error {
	ti.tableWriterBase.init(txn, ti.tableDesc(), evalCtx, sv)
	return nil
}

// row is part of the tableWriter interface.
func (ti *tableInserter) row(
	ctx context.Context, values tree.Datums, pm row.PartialIndexUpdateHelper, traceKV bool,
) error {
	ti.currentBatchSize++
	return ti.ri.InsertRow(ctx, ti.b, values, pm, false /* overwrite */, traceKV)
}

// tableDesc is part of the tableWriter interface.
func (ti *tableInserter) tableDesc() catalog.TableDescriptor {
	return ti.ri.Helper.TableDesc
}

// walkExprs is part of the tableWriter interface.
func (ti *tableInserter) walkExprs(_ func(desc string, index int, expr tree.TypedExpr)) {}
