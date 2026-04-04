// Copyright 2019 The Cockroach Authors.
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

package optbuilder

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

func (b *Builder) buildControlJobs(n *tree.ControlJobs, inScope *scope) (outScope *scope) {
	// We don't allow the input statement to reference outer columns, so we
	// pass a "blank" scope rather than inScope.
	emptyScope := b.allocScope()
	colTypes := []*types.T{types.Int}
	inputScope := b.buildStmt(n.Jobs, colTypes, emptyScope)

	var reason opt.ScalarExpr
	if n.Reason != nil {
		reasonStr := emptyScope.resolveType(n.Reason, types.String)
		reason = b.buildScalar(
			reasonStr, emptyScope, nil /* outScope */, nil /* outCol */, nil, /* colRefs */
		)
	} else {
		reason = b.factory.ConstructNull(types.String)
	}

	checkInputColumns(
		fmt.Sprintf("%s JOBS", tree.JobCommandToStatement[n.Command]),
		inputScope,
		[]string{"job_id"},
		colTypes,
		1, /* minPrefix */
	)
	outScope = inScope.push()
	outScope.expr = b.factory.ConstructControlJobs(
		inputScope.expr,
		reason,
		&memo.ControlJobsPrivate{
			Props:   inputScope.makePhysicalProps(),
			Command: n.Command,
		},
	)
	return outScope
}

func (b *Builder) buildCancelQueries(n *tree.CancelQueries, inScope *scope) (outScope *scope) {
	// We don't allow the input statement to reference outer columns, so we
	// pass a "blank" scope rather than inScope.
	emptyScope := b.allocScope()
	colTypes := []*types.T{types.String}
	inputScope := b.buildStmt(n.Queries, colTypes, emptyScope)

	checkInputColumns(
		"CANCEL QUERIES",
		inputScope,
		[]string{"query_id"},
		colTypes,
		1, /* minPrefix */
	)
	outScope = inScope.push()
	outScope.expr = b.factory.ConstructCancelQueries(
		inputScope.expr,
		&memo.CancelPrivate{
			Props:    inputScope.makePhysicalProps(),
			IfExists: n.IfExists,
		},
	)
	return outScope
}

func (b *Builder) buildCancelSessions(n *tree.CancelSessions, inScope *scope) (outScope *scope) {
	// We don't allow the input statement to reference outer columns, so we
	// pass a "blank" scope rather than inScope.
	emptyScope := b.allocScope()
	colTypes := []*types.T{types.String}
	inputScope := b.buildStmt(n.Sessions, colTypes, emptyScope)

	checkInputColumns(
		"CANCEL SESSIONS",
		inputScope,
		[]string{"session_id"},
		colTypes,
		1, /* minPrefix */
	)
	outScope = inScope.push()
	outScope.expr = b.factory.ConstructCancelSessions(
		inputScope.expr,
		&memo.CancelPrivate{
			Props:    inputScope.makePhysicalProps(),
			IfExists: n.IfExists,
		},
	)
	return outScope
}

func (b *Builder) buildControlSchedules(
	n *tree.ControlSchedules, inScope *scope,
) (outScope *scope) {
	if err := b.catalog.RequireAdminRole(b.ctx, n.StatementTag()); err != nil {
		panic(err)
	}

	// We don't allow the input statement to reference outer columns, so we
	// pass a "blank" scope rather than inScope.
	emptyScope := b.allocScope()
	colTypes := []*types.T{types.Int}
	inputScope := b.buildStmt(n.Schedules, colTypes, emptyScope)

	checkInputColumns(
		fmt.Sprintf("%s SCHEDULES", n.Command),
		inputScope,
		[]string{"schedule_id"},
		colTypes,
		1, /* minPrefix */
	)

	outScope = inScope.push()
	outScope.expr = b.factory.ConstructControlSchedules(
		inputScope.expr,
		&memo.ControlSchedulesPrivate{
			Props:   inputScope.makePhysicalProps(),
			Command: n.Command,
		},
	)
	return outScope
}

func (b *Builder) buildCreateStatistics(n *tree.CreateStats, inScope *scope) (outScope *scope) {
	outScope = inScope.push()
	outScope.expr = b.factory.ConstructCreateStatistics(&memo.CreateStatisticsPrivate{
		Syntax: n,
	})
	return outScope
}
