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
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/screl"
	"github.com/cockroachdb/cockroach/pkg/util/iterutil"
)

// NodeIterator is used to iterate nodes. Return iterutil.StopIteration to
// return early with no error.
type NodeIterator func(n *screl.Node) error

// ForEachNode iterates the nodes in the graph.
func (g *Graph) ForEachNode(it NodeIterator) error {
	for _, m := range g.targetNodes {
		for i := 0; i < scpb.NumStatus; i++ {
			if ts, ok := m[scpb.Status(i)]; ok {
				if err := it(ts); err != nil {
					if iterutil.Done(err) {
						err = nil
					}
					return err
				}
			}
		}
	}
	return nil
}

// EdgeIterator is used to iterate edges. Return iterutil.StopIteration to
// return early with no error.
type EdgeIterator func(e Edge) error

// ForEachEdge iterates the edges in the graph.
func (g *Graph) ForEachEdge(it EdgeIterator) error {
	for _, e := range g.edges {
		if err := it(e); err != nil {
			if iterutil.Done(err) {
				err = nil
			}
			return err
		}
	}
	return nil
}

// DepEdgeIterator is used to iterate dep edges. Return iterutil.StopIteration
// to return early with no error.
type DepEdgeIterator func(de *DepEdge) error

// ForEachDepEdgeFrom iterates the dep edges in the graph with the selected
// source.
func (g *Graph) ForEachDepEdgeFrom(n *screl.Node, it DepEdgeIterator) (err error) {
	return g.depEdgesFrom.iterateSourceNode(n, it)
}

// ForEachDepEdgeTo iterates the dep edges in the graph with the selected
// destination.
func (g *Graph) ForEachDepEdgeTo(n *screl.Node, it DepEdgeIterator) (err error) {
	return g.depEdgesTo.iterateSourceNode(n, it)
}
