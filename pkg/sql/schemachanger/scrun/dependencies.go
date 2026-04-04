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

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scexec"
)

// JobTxnFunc is used to run a transactional stage of a schema change on
// behalf of a job. See JobRunDependencies.WithTxnInJob().
type JobTxnFunc = func(ctx context.Context, txnDeps scexec.Dependencies) error

// JobRunDependencies contains the dependencies required for
// executing the schema change job, i.e. for the logic in its Resume() method.
type JobRunDependencies interface {
	// WithTxnInJob is a wrapper for opening and committing a transaction around
	// the execution of the callback. After committing the transaction, the job
	// registry should be notified to adopt jobs.
	WithTxnInJob(ctx context.Context, fn JobTxnFunc) error
}
