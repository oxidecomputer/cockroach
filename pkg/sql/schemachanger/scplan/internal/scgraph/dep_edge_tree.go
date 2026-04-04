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

package scgraph

import (
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/screl"
	"github.com/cockroachdb/cockroach/pkg/util/iterutil"
	"github.com/google/btree"
)

type depEdgeTree struct {
	t     *btree.BTree
	order edgeTreeOrder
	cmp   nodeCmpFn
}

type nodeCmpFn func(a, b *screl.Node) (less, eq bool)

func newDepEdgeTree(order edgeTreeOrder, cmp nodeCmpFn) *depEdgeTree {
	const degree = 8 // arbitrary
	return &depEdgeTree{
		t:     btree.New(degree),
		order: order,
		cmp:   cmp,
	}
}

// edgeTreeOrder order in which the edge tree is sorted,
// either based on from/to node indexes.
type edgeTreeOrder bool

func (o edgeTreeOrder) first(e Edge) *screl.Node {
	if o == fromTo {
		return e.From()
	}
	return e.To()
}

func (o edgeTreeOrder) second(e Edge) *screl.Node {
	if o == toFrom {
		return e.From()
	}
	return e.To()
}

const (
	fromTo edgeTreeOrder = true
	toFrom edgeTreeOrder = false
)

// edgeTreeEntry BTree items for tracking edges
// in an ordered manner.
type edgeTreeEntry struct {
	t    *depEdgeTree
	edge *DepEdge
}

func (et *depEdgeTree) insert(e *DepEdge) {
	et.t.ReplaceOrInsert(&edgeTreeEntry{
		t:    et,
		edge: e,
	})
}

func (et *depEdgeTree) iterateSourceNode(n *screl.Node, it DepEdgeIterator) (err error) {
	e := &edgeTreeEntry{t: et, edge: &DepEdge{}}
	if et.order == fromTo {
		e.edge.from = n
	} else {
		e.edge.to = n
	}
	et.t.AscendGreaterOrEqual(e, func(i btree.Item) (wantMore bool) {
		e := i.(*edgeTreeEntry)
		if et.order.first(e.edge) != n {
			return false
		}
		err = it(e.edge)
		return err == nil
	})
	if iterutil.Done(err) {
		err = nil
	}
	return err
}

// Less implements btree.Item.
func (e *edgeTreeEntry) Less(otherItem btree.Item) bool {
	o := otherItem.(*edgeTreeEntry)
	if less, eq := e.t.cmp(e.t.order.first(e.edge), e.t.order.first(o.edge)); !eq {
		return less
	}
	less, _ := e.t.cmp(e.t.order.second(e.edge), e.t.order.second(o.edge))
	return less
}
