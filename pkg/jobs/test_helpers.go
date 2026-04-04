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

package jobs

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/server/tracedumper"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// TestingNudgeAdoptionQueue is used by tests to tell the registry that there is
// a job to be adopted.
func (r *Registry) TestingNudgeAdoptionQueue() {
	r.adoptionCh <- claimAndResumeClaimedJobs
}

type config struct {
	jobID jobspb.JobID
}

// TestCreateAndStartJobOption optionally modifies TestingCreateAndStartJob.
type TestCreateAndStartJobOption func(*config)

// WithJobID is used to inject an existing JobID to TestingCreateAndStartJob.
func WithJobID(jobID jobspb.JobID) TestCreateAndStartJobOption {
	return func(c *config) {
		c.jobID = jobID
	}
}

// TestingCreateAndStartJob creates and asynchronously starts a job from record.
// An error is returned if the job type has not been registered with
// RegisterConstructor. The ctx passed to this function is not the context the
// job will be started with (canceling ctx will not cause the job to cancel).
func TestingCreateAndStartJob(
	ctx context.Context, r *Registry, db *kv.DB, record Record, opts ...TestCreateAndStartJobOption,
) (*StartableJob, error) {
	var rj *StartableJob
	c := config{
		jobID: r.MakeJobID(),
	}
	for _, opt := range opts {
		opt(&c)
	}
	if err := db.Txn(ctx, func(ctx context.Context, txn *kv.Txn) (err error) {
		return r.CreateStartableJobWithTxn(ctx, &rj, c.jobID, txn, record)
	}); err != nil {
		if rj != nil {
			if cleanupErr := rj.CleanupOnRollback(ctx); cleanupErr != nil {
				log.Warningf(ctx, "failed to cleanup StartableJob: %v", cleanupErr)
			}
		}
		return nil, err
	}
	err := rj.Start(ctx)
	if err != nil {
		return nil, err
	}
	return rj, nil
}

// TestingGetTraceDumpDir returns the directory in which jobs might dump their
// traces after execution.
func TestingGetTraceDumpDir(r *Registry) string {
	if r.td == nil {
		return ""
	}
	return tracedumper.TestingGetTraceDumpDir(r.td)
}
