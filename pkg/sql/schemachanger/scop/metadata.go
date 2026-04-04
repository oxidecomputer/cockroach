// Copyright 2022 The Cockroach Authors.
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

// StatementForDropJob is a statement used to build a description for a
// drop job. The set of statements associated with the drop job will
// be accumulated for the description.
type StatementForDropJob struct {

	// Statement is the statement which lead to this drop.
	Statement string

	// StatementID is the order of the statement in the transaction. It is used
	// to synthesize the appropriate description.
	StatementID uint32

	// Rollback should be marked true if the schema change job is currently
	// rolling back. This is needed to build the correct description for the
	// job.
	Rollback bool
}
