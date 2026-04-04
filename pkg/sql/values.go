// Copyright 2015 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/sql/rowcontainer"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

type valuesNode struct {
	// Note: the columns can be renamed in place (see planMutableColumns).
	columns colinfo.ResultColumns
	tuples  [][]tree.TypedExpr

	// specifiedInQuery is set if the valuesNode represents a literal
	// relational expression that was present in the original SQL text,
	// as opposed to e.g. a valuesNode resulting from the expansion of
	// a vtable value generator. This changes distsql physical planning.
	specifiedInQuery bool

	valuesRun
}

func (p *planner) newContainerValuesNode(columns colinfo.ResultColumns, capacity int) *valuesNode {
	return &valuesNode{
		columns: columns,
		valuesRun: valuesRun{
			rows: rowcontainer.NewRowContainerWithCapacity(
				p.EvalContext().Mon.MakeBoundAccount(),
				colinfo.ColTypeInfoFromResCols(columns),
				capacity,
			),
		},
	}
}

// valuesRun is the run-time state of a valuesNode during local execution.
type valuesRun struct {
	rows    *rowcontainer.RowContainer
	nextRow int // The index of the next row.
}

func (n *valuesNode) startExec(params runParams) error {
	if n.rows != nil {
		// n.rows was already created in newContainerValuesNode.
		// Nothing to do here.
		return nil
	}

	// This node is coming from a SQL query (as opposed to sortNode and
	// others that create a valuesNode internally for storing results
	// from other planNodes), so its expressions need evaluating.
	// This may run subqueries.
	n.rows = rowcontainer.NewRowContainerWithCapacity(
		params.extendedEvalCtx.Mon.MakeBoundAccount(),
		colinfo.ColTypeInfoFromResCols(n.columns),
		len(n.tuples),
	)

	row := make([]tree.Datum, len(n.columns))
	for _, tupleRow := range n.tuples {
		for i, typedExpr := range tupleRow {
			var err error
			row[i], err = typedExpr.Eval(params.EvalContext())
			if err != nil {
				return err
			}
		}
		if _, err := n.rows.AddRow(params.ctx, row); err != nil {
			return err
		}
	}

	return nil
}

func (n *valuesNode) Next(runParams) (bool, error) {
	if n.nextRow >= n.rows.Len() {
		return false, nil
	}
	n.nextRow++
	return true, nil
}

func (n *valuesNode) Values() tree.Datums {
	return n.rows.At(n.nextRow - 1)
}

func (n *valuesNode) Close(ctx context.Context) {
	if n.rows != nil {
		n.rows.Close(ctx)
		n.rows = nil
	}
}
