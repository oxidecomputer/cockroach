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

package scop

// A Phase represents the context in which an op is executed within a schema
// change. Different phases require different dependencies for the execution of
// the ops to be plumbed in.
//
// Today, we support the phases corresponding to async schema changes initiated
// and partially executed in the user transaction. This will change as we
// transition to transactional schema changes.
type Phase int

//go:generate stringer --type Phase

const (
	_ Phase = iota
	// StatementPhase refers to execution of ops occurring during statement
	// execution during the user transaction.
	StatementPhase
	// PreCommitPhase refers to execution of ops occurring during the user
	// transaction immediately before commit.
	PreCommitPhase
	// PostCommitPhase refers to execution of ops occurring after the user
	// transaction has committed (i.e., in the async schema change job).
	// Note: Planning rules cannot ever be in this phase, since all those operations
	// should be executed in pre-commit.
	PostCommitPhase
	// PostCommitNonRevertiblePhase is like PostCommitPhase but in which target
	// status changes are non-revertible.
	PostCommitNonRevertiblePhase

	// EarliestPhase references the earliest possible execution phase.
	EarliestPhase = StatementPhase

	// LatestPhase references the latest possible execution phase.
	LatestPhase = PostCommitNonRevertiblePhase
)
