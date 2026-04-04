// Copyright 2020 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

type invertedJoinNode struct {
	input planNode
	table *scanNode

	// joinType is one of INNER, LEFT_OUTER, LEFT_SEMI, LEFT_ANTI.
	joinType descpb.JoinType

	// prefixEqCols identifies the columns from the input which are used for the
	// lookup if the index is a multi-column inverted index. These correspond to
	// the non-inverted prefix columns of the index we are looking up. This is
	// empty if the index is not a multi-column inverted index.
	prefixEqCols []int

	// The inverted expression to evaluate.
	invertedExpr tree.TypedExpr

	// columns are the produced columns, namely the input columns and (unless the
	// join type is semi or anti join) the columns in the table scanNode. It can
	// include an additional continuation column for paired joins.
	columns colinfo.ResultColumns

	// onExpr is any ON condition to be used in conjunction with the inverted
	// expression.
	onExpr tree.TypedExpr

	isFirstJoinInPairedJoiner bool

	reqOrdering ReqOrdering
}

func (ij *invertedJoinNode) startExec(params runParams) error {
	panic("invertedJoinNode cannot be run in local mode")
}

func (ij *invertedJoinNode) Next(params runParams) (bool, error) {
	panic("invertedJoinNode cannot be run in local mode")
}

func (ij *invertedJoinNode) Values() tree.Datums {
	panic("invertedJoinNode cannot be run in local mode")
}

func (ij *invertedJoinNode) Close(ctx context.Context) {
	ij.input.Close(ctx)
	ij.table.Close(ctx)
}
