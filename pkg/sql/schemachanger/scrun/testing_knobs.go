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

package scrun

import "github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scplan"

// TestingKnobs are testing knobs which affect the running of declarative
// schema changes.
type TestingKnobs struct {
	// BeforeStage is called before ops passed to the executor are executed.
	// Errors returned are injected into the executor.
	BeforeStage func(p scplan.Plan, stageIdx int) error

	// AfterStage is invoked after all ops are executed.
	// Errors returned are injected into the executor.
	AfterStage func(p scplan.Plan, stageIdx int) error

	// BeforeWaitingForConcurrentSchemaChanges is called at the start of waiting
	// for concurrent schema changes to finish.
	BeforeWaitingForConcurrentSchemaChanges func(stmts []string)

	// OnPostCommitError is called whenever the schema changer job returns an
	// error.
	OnPostCommitError func(p scplan.Plan, stageIdx int, err error) error
}

// ModuleTestingKnobs is part of the base.ModuleTestingKnobs interface.
func (*TestingKnobs) ModuleTestingKnobs() {}
