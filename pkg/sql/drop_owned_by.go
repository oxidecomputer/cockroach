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

	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
	"github.com/cockroachdb/cockroach/pkg/util/errorutil/unimplemented"
)

// dropOwnedByNode represents a DROP OWNED BY <role(s)> statement.
type dropOwnedByNode struct {
	// TODO(angelaw): Uncomment when implementing - commenting out due to linting error.
	//n *tree.DropOwnedBy
}

func (p *planner) DropOwnedBy(ctx context.Context) (planNode, error) {
	if err := checkSchemaChangeEnabled(
		ctx,
		p.ExecCfg(),
		"DROP OWNED BY",
	); err != nil {
		return nil, err
	}
	telemetry.Inc(sqltelemetry.CreateDropOwnedByCounter())
	// TODO(angelaw): Implementation.
	return nil, unimplemented.NewWithIssue(55381, "drop owned by is not yet implemented")
}

func (n *dropOwnedByNode) startExec(params runParams) error {
	// TODO(angelaw): Implementation.
	return nil
}
func (n *dropOwnedByNode) Next(runParams) (bool, error) { return false, nil }
func (n *dropOwnedByNode) Values() tree.Datums          { return tree.Datums{} }
func (n *dropOwnedByNode) Close(context.Context)        {}
