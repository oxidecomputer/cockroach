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

package gcjob

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
)

// gcTenant drops the data of tenant that has an expired deadline and updates
// the job details to mark the work it did. The job progress is updated in
// place, but needs to be persisted to the job.
func gcTenant(
	ctx context.Context,
	execCfg *sql.ExecutorConfig,
	tenID uint64,
	progress *jobspb.SchemaChangeGCProgress,
) error {
	if log.V(2) {
		log.Infof(ctx, "GC is being considered for tenant: %d", tenID)
	}

	if progress.Tenant.Status == jobspb.SchemaChangeGCProgress_WAITING_FOR_GC {
		return errors.AssertionFailedf(
			"Tenant id %d is expired and should not be in state %+v",
			tenID, jobspb.SchemaChangeGCProgress_WAITING_FOR_GC,
		)
	}

	info, err := sql.GetTenantRecord(ctx, execCfg, nil /* txn */, tenID)
	if err != nil {
		if pgerror.GetPGCode(err) == pgcode.UndefinedObject {
			// The tenant row is deleted only after its data is cleared so there is
			// nothing to do in this case but mark the job as done.
			if progress.Tenant.Status != jobspb.SchemaChangeGCProgress_DELETED {
				// This will happen if the job deletes the tenant row and fails to update
				// its progress. In this case there's nothing to do but update the job
				// progress.
				log.Errorf(ctx, "tenant id %d not found while attempting to GC", tenID)
				progress.Tenant.Status = jobspb.SchemaChangeGCProgress_DELETED
			}
			return nil
		}
		return errors.Wrapf(err, "fetching tenant %d", info.ID)
	}

	// This case should never happen.
	if progress.Tenant.Status == jobspb.SchemaChangeGCProgress_DELETED {
		return errors.AssertionFailedf("GC state for tenant %+v is DELETED yet the tenant row still exists", info)
	}

	if err := sql.GCTenantSync(ctx, execCfg, info); err != nil {
		return errors.Wrapf(err, "gc tenant %d", info.ID)
	}

	progress.Tenant.Status = jobspb.SchemaChangeGCProgress_DELETED
	return nil
}
