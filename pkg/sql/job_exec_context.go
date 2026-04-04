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
	"github.com/cockroachdb/cockroach/pkg/migration"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/lease"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondatapb"
)

// plannerJobExecContext is a wrapper to implement JobExecContext with a planner
// without allowing casting directly to a planner. Eventually it would be nice
// if we could implement the API entirely without a planner however the only
// implementation of extendedEvalContext is very tied to a planner.
type plannerJobExecContext struct {
	p *planner
}

// MakeJobExecContext makes a JobExecContext.
func MakeJobExecContext(
	opName string, user security.SQLUsername, memMetrics *MemoryMetrics, execCfg *ExecutorConfig,
) (JobExecContext, func()) {
	plannerInterface, close := NewInternalPlanner(
		opName,
		nil, /*txn*/
		user,
		memMetrics,
		execCfg,
		sessiondatapb.SessionData{},
	)
	p := plannerInterface.(*planner)
	return &plannerJobExecContext{p: p}, close
}

func (e *plannerJobExecContext) SemaCtx() *tree.SemaContext { return e.p.SemaCtx() }
func (e *plannerJobExecContext) ExtendedEvalContext() *extendedEvalContext {
	return e.p.ExtendedEvalContext()
}
func (e *plannerJobExecContext) SessionData() *sessiondata.SessionData {
	return e.p.SessionData()
}
func (e *plannerJobExecContext) SessionDataMutatorIterator() *sessionDataMutatorIterator {
	return e.p.SessionDataMutatorIterator()
}
func (e *plannerJobExecContext) ExecCfg() *ExecutorConfig        { return e.p.ExecCfg() }
func (e *plannerJobExecContext) DistSQLPlanner() *DistSQLPlanner { return e.p.DistSQLPlanner() }
func (e *plannerJobExecContext) LeaseMgr() *lease.Manager        { return e.p.LeaseMgr() }
func (e *plannerJobExecContext) User() security.SQLUsername      { return e.p.User() }
func (e *plannerJobExecContext) MigrationJobDeps() migration.JobDeps {
	return e.p.MigrationJobDeps()
}
func (e *plannerJobExecContext) SpanConfigReconciler() spanconfig.Reconciler {
	return e.p.SpanConfigReconciler()
}

// JobExecContext provides the execution environment for a job. It is what is
// passed to the Resume/OnFailOrCancel/OnPauseRequested methods of a jobs's
// Resumer to give that resumer access to things like ExecutorCfg, LeaseMgr,
// etc -- the kinds of things that would usually be on planner or similar during
// a non-job SQL statement's execution. Unlike a planner however, or planner-ish
// interfaces like PlanHookState, JobExecContext does not include a txn or the
// methods that defined in terms of "the" txn, such as privilege/name accessors.
// (though note that ExtendedEvalContext may transitively include methods that
// close over/expect a txn so use it with caution).
type JobExecContext interface {
	SemaCtx() *tree.SemaContext
	ExtendedEvalContext() *extendedEvalContext
	SessionData() *sessiondata.SessionData
	SessionDataMutatorIterator() *sessionDataMutatorIterator
	ExecCfg() *ExecutorConfig
	DistSQLPlanner() *DistSQLPlanner
	LeaseMgr() *lease.Manager
	User() security.SQLUsername
	MigrationJobDeps() migration.JobDeps
	SpanConfigReconciler() spanconfig.Reconciler
}
