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
	"strings"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// showVarNode represents a SHOW <var> statement.
// This is reached if <var> contains a period.
type showVarNode struct {
	name  string
	shown bool
	val   string
}

func (s *showVarNode) startExec(params runParams) error {
	return nil
}

func (s *showVarNode) Next(params runParams) (bool, error) {
	if s.shown {
		return false, nil
	}
	s.shown = true

	_, v, err := getSessionVar(s.name, false /* missingOk */)
	if err != nil {
		return false, err
	}
	s.val, err = v.Get(params.extendedEvalCtx)
	return true, err
}

func (s *showVarNode) Values() tree.Datums {
	return tree.Datums{tree.NewDString(s.val)}
}

func (s *showVarNode) Close(ctx context.Context) {}

// ShowVar shows a session variable.
func (p *planner) ShowVar(ctx context.Context, n *tree.ShowVar) (planNode, error) {
	return &showVarNode{name: strings.ToLower(n.Name)}, nil
}
