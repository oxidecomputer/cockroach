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
	"github.com/cockroachdb/cockroach/pkg/migration"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/lease"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
)

// FakeJobExecContext is used for mocking the JobExecContext in tests.
type FakeJobExecContext struct {
	JobExecContext
	ExecutorConfig *ExecutorConfig
}

// ExecCfg implements the JobExecContext interface.
func (p *FakeJobExecContext) ExecCfg() *ExecutorConfig {
	return p.ExecutorConfig
}

// SemaCtx implements the JobExecContext interface.
func (p *FakeJobExecContext) SemaCtx() *tree.SemaContext {
	return nil
}

// ExtendedEvalContext implements the JobExecContext interface.
func (p *FakeJobExecContext) ExtendedEvalContext() *extendedEvalContext {
	panic("unimplemented")
}

// SessionData implements the JobExecContext interface.
func (p *FakeJobExecContext) SessionData() *sessiondata.SessionData {
	return nil
}

// SessionDataMutatorIterator implements the JobExecContext interface.
func (p *FakeJobExecContext) SessionDataMutatorIterator() *sessionDataMutatorIterator {
	panic("unimplemented")
}

// DistSQLPlanner implements the JobExecContext interface.
func (p *FakeJobExecContext) DistSQLPlanner() *DistSQLPlanner {
	panic("unimplemented")
}

// LeaseMgr implements the JobExecContext interface.
func (p *FakeJobExecContext) LeaseMgr() *lease.Manager {
	panic("unimplemented")
}

// User implements the JobExecContext interface.
func (p *FakeJobExecContext) User() security.SQLUsername {
	panic("unimplemented")
}

// MigrationJobDeps implements the JobExecContext interface.
func (p *FakeJobExecContext) MigrationJobDeps() migration.JobDeps {
	panic("unimplemented")
}

// SpanConfigReconciler implements the JobExecContext interface.
func (p *FakeJobExecContext) SpanConfigReconciler() spanconfig.Reconciler {
	panic("unimplemented")
}
