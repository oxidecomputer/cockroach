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

package opgen

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan/internal/scgraph"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/screl"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

type registry struct {
	targets []target
}

var opRegistry = &registry{}

// BuildGraph constructs a graph with operation edges populated from an initial
// state.
func BuildGraph(cs scpb.CurrentState) (*scgraph.Graph, error) {
	return opRegistry.buildGraph(cs)
}

func (r *registry) buildGraph(cs scpb.CurrentState) (_ *scgraph.Graph, err error) {
	start := timeutil.Now()
	defer func() {
		if err != nil || !log.V(2) {
			return
		}
		log.Infof(context.TODO(), "operation graph generation took %v", timeutil.Since(start))
	}()
	g, err := scgraph.New(cs)
	if err != nil {
		return nil, err
	}
	// Iterate through each match of initial state target's to target rules
	// and apply the relevant op edges to the graph. Copy out the elements
	// to not mutate the database in place.
	type toAdd struct {
		transition
		n *screl.Node
	}
	var edgesToAdd []toAdd
	md := makeTargetsWithElementMap(cs)
	for _, t := range r.targets {
		edgesToAdd = edgesToAdd[:0]
		if err := t.iterateFunc(g.Database(), func(n *screl.Node) error {
			status := n.CurrentStatus
			for _, op := range t.transitions {
				if op.from == status {
					edgesToAdd = append(edgesToAdd, toAdd{
						transition: op,
						n:          n,
					})
					status = op.to
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
		for _, e := range edgesToAdd {
			var ops []scop.Op
			if e.ops != nil {
				ops = e.ops(e.n.Element(), md)
			}
			if err := g.AddOpEdges(
				e.n.Target, e.from, e.to, e.revertible, e.minPhase, ops...,
			); err != nil {
				return nil, err
			}
		}

	}
	return g, nil
}
