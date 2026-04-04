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
	"github.com/cockroachdb/cockroach/pkg/sql/inverted"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

type invertedFilterNode struct {
	input           planNode
	expression      *inverted.SpanExpression
	preFiltererExpr tree.TypedExpr
	preFiltererType *types.T
	invColumn       int
	resultColumns   colinfo.ResultColumns
}

func (n *invertedFilterNode) startExec(params runParams) error {
	panic("invertedFiltererNode can't be run in local mode")
}
func (n *invertedFilterNode) Close(ctx context.Context) {
	n.input.Close(ctx)
}
func (n *invertedFilterNode) Next(params runParams) (bool, error) {
	panic("invertedFiltererNode can't be run in local mode")
}
func (n *invertedFilterNode) Values() tree.Datums {
	panic("invertedFiltererNode can't be run in local mode")
}
