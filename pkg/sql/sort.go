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
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// sortNode represents a node that sorts the rows returned by its
// sub-node.
type sortNode struct {
	plan     planNode
	ordering colinfo.ColumnOrdering
	// When alreadyOrderedPrefix is non-zero, the input is already ordered on
	// the prefix ordering[:alreadyOrderedPrefix].
	alreadyOrderedPrefix int
}

func (n *sortNode) startExec(runParams) error {
	panic("sortNode cannot be run in local mode")
}

func (n *sortNode) Next(params runParams) (bool, error) {
	panic("sortNode cannot be run in local mode")
}

func (n *sortNode) Values() tree.Datums {
	panic("sortNode cannot be run in local mode")
}

func (n *sortNode) Close(ctx context.Context) {
	n.plan.Close(ctx)
}
