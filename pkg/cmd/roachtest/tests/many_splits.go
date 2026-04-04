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

package tests

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
	"github.com/stretchr/testify/require"
)

// runManySplits attempts to create 2000 tiny ranges on a 4-node cluster using
// left-to-right splits and check the cluster is still live afterwards.
func runManySplits(ctx context.Context, t test.Test, c cluster.Cluster) {
	c.Put(ctx, t.Cockroach(), "./cockroach")
	settings := install.MakeClusterSettings()
	settings.Env = append(settings.Env, "COCKROACH_SCAN_MAX_IDLE_TIME=5ms")
	c.Start(ctx, t.L(), option.DefaultStartOpts(), settings)

	db := c.Conn(ctx, t.L(), 1)
	defer db.Close()

	// Wait for up-replication then create many ranges.
	err := WaitFor3XReplication(ctx, t, db)
	require.NoError(t, err)

	m := c.NewMonitor(ctx, c.All())
	m.Go(func(ctx context.Context) error {
		const numRanges = 2000
		t.L().Printf("creating %d ranges...", numRanges)
		if _, err := db.ExecContext(ctx, fmt.Sprintf(`
			CREATE TABLE t(x, PRIMARY KEY(x)) AS TABLE generate_series(1,%[1]d);
            ALTER TABLE t SPLIT AT TABLE generate_series(1,%[1]d);
		`, numRanges)); err != nil {
			return err
		}
		return nil
	})
	m.Wait()
}
