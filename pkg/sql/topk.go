// Copyright 2021 The Cockroach Authors.
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

// topKNode represents a node that returns only the top K rows according to the
// ordering, in the order specified.
type topKNode struct {
	plan     planNode
	k        int64
	ordering colinfo.ColumnOrdering
	// When alreadyOrderedPrefix is non-zero, the input is already ordered on
	// the prefix ordering[:alreadyOrderedPrefix].
	alreadyOrderedPrefix int
}

func (n *topKNode) startExec(params runParams) error {
	panic("topKNode cannot be run in local mode")
}

func (n *topKNode) Next(params runParams) (bool, error) {
	panic("topKNode cannot be run in local mode")
}

func (n *topKNode) Values() tree.Datums {
	panic("topKNode cannot be run in local mode")
}

func (n *topKNode) Close(ctx context.Context) {
	n.plan.Close(ctx)
}
