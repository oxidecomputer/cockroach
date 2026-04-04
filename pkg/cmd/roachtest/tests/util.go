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

package tests

import (
	"context"
	gosql "database/sql"
	"fmt"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/util/retry"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
)

// WaitFor3XReplication is like WaitForReplication but specifically requires
// three as the minimum number of voters a range must be replicated on.
func WaitFor3XReplication(ctx context.Context, t test.Test, db *gosql.DB) error {
	return WaitForReplication(ctx, t, db, 3 /* replicationFactor */)
}

// WaitForReplication waits until all ranges in the system are on at least
// replicationFactor voters.
func WaitForReplication(
	ctx context.Context, t test.Test, db *gosql.DB, replicationFactor int,
) error {
	t.L().Printf("waiting for initial up-replication...")
	tStart := timeutil.Now()
	var oldN int
	for {
		var n int
		if err := db.QueryRowContext(
			ctx,
			fmt.Sprintf(
				"SELECT count(1) FROM crdb_internal.ranges WHERE array_length(replicas, 1) < %d",
				replicationFactor,
			),
		).Scan(&n); err != nil {
			return err
		}
		if n == 0 {
			t.L().Printf("up-replication complete")
			return nil
		}
		if timeutil.Since(tStart) > 30*time.Second || oldN != n {
			t.L().Printf("still waiting for full replication (%d ranges left)", n)
		}
		oldN = n
		time.Sleep(time.Second)
	}
}

// WaitForUpdatedReplicationReport waits for an updated replication report.
func WaitForUpdatedReplicationReport(ctx context.Context, t test.Test, db *gosql.DB) {
	t.L().Printf("waiting for updated replication report...")

	// Temporarily drop the replication report interval down.
	if _, err := db.ExecContext(
		ctx, `SET CLUSTER setting kv.replication_reports.interval = '2s'`,
	); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := db.ExecContext(
			ctx, `RESET CLUSTER setting kv.replication_reports.interval`,
		); err != nil {
			t.Fatal(err)
		}
	}()

	// Wait for a new report with a timestamp after tStart to ensure
	// that the report picks up any new tables or zones.
	tStart := timeutil.Now()
	for r := retry.StartWithCtx(ctx, retry.Options{}); r.Next(); {
		var count int
		var gen gosql.NullTime
		if err := db.QueryRowContext(
			ctx, `SELECT count(*), min(generated) FROM system.reports_meta`,
		).Scan(&count, &gen); err != nil {
			if !errors.Is(err, gosql.ErrNoRows) {
				t.Fatal(err)
			}
			// No report generated yet. There are 3 types of reports. We want to
			// see a result for all of them.
		} else if count == 3 && tStart.Before(gen.Time) {
			// New report generated.
			return
		}
		if timeutil.Since(tStart) > 30*time.Second {
			t.L().Printf("still waiting for updated replication report")
		}
	}
}

// SetAdmissionControl sets the admission control cluster settings on the
// given cluster.
func SetAdmissionControl(ctx context.Context, t test.Test, c cluster.Cluster, enabled bool) {
	db := c.Conn(ctx, t.L(), 1)
	defer db.Close()
	val := "true"
	if !enabled {
		val = "false"
	}
	for _, setting := range []string{"admission.kv.enabled", "admission.sql_kv_response.enabled",
		"admission.sql_sql_response.enabled"} {
		if _, err := db.ExecContext(
			ctx, "SET CLUSTER SETTING "+setting+" = '"+val+"'"); err != nil {
			t.Fatalf("failed to set admission control to %t: %v", enabled, err)
		}
	}
}
