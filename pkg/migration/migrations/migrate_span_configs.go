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

package migrations

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/migration"
	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/retry"
	"github.com/cockroachdb/errors"
)

func ensureSpanConfigReconciliation(
	ctx context.Context, _ clusterversion.ClusterVersion, d migration.TenantDeps, j *jobs.Job,
) error {
	if !d.Codec.ForSystemTenant() {
		return nil
	}

	retryOpts := retry.Options{
		InitialBackoff: time.Second,
		MaxBackoff:     2 * time.Second,
		Multiplier:     2,
		MaxRetries:     50,
	}
	for r := retry.StartWithCtx(ctx, retryOpts); r.Next(); {
		row, err := d.InternalExecutor.QueryRowEx(ctx, "get-spanconfig-progress", nil,
			sessiondata.NodeUserSessionDataOverride, `
SELECT progress
  FROM system.jobs
 WHERE id = (SELECT job_id FROM [SHOW AUTOMATIC JOBS] WHERE job_type = 'AUTO SPAN CONFIG RECONCILIATION')
`)
		if err != nil {
			return err
		}
		if row == nil {
			log.Info(ctx, "reconciliation job not found, retrying...")
			continue
		}
		progress, err := jobs.UnmarshalProgress(row[0])
		if err != nil {
			return err
		}
		sp, ok := progress.GetDetails().(*jobspb.Progress_AutoSpanConfigReconciliation)
		if !ok {
			log.Fatal(ctx, "unexpected job progress type")
		}
		if sp.AutoSpanConfigReconciliation.Checkpoint.IsEmpty() {
			log.Info(ctx, "waiting for span config reconciliation...")
			continue
		}

		return nil
	}

	return errors.Newf("unable to reconcile span configs")
}

func ensureSpanConfigSubscription(
	ctx context.Context, _ clusterversion.ClusterVersion, deps migration.SystemDeps, _ *jobs.Job,
) error {
	return deps.Cluster.UntilClusterStable(ctx, func() error {
		return deps.Cluster.ForEveryNode(ctx, "ensure-span-config-subscription",
			func(ctx context.Context, client serverpb.MigrationClient) error {
				req := &serverpb.WaitForSpanConfigSubscriptionRequest{}
				_, err := client.WaitForSpanConfigSubscription(ctx, req)
				return err
			})
	})
}
