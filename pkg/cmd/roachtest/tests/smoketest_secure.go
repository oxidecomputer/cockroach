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

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/stretchr/testify/require"
)

func registerSecure(r registry.Registry) {
	for _, numNodes := range []int{1, 3} {
		r.Add(registry.TestSpec{
			Name:    fmt.Sprintf("smoketest/secure/nodes=%d", numNodes),
			Tags:    []string{"smoketest", "weekly"},
			Owner:   registry.OwnerKV, // TODO: OwnerTestEng once the open PR that introduces it has merged
			Cluster: r.MakeClusterSpec(numNodes),
			Run: func(ctx context.Context, t test.Test, c cluster.Cluster) {
				c.Put(ctx, t.Cockroach(), "./cockroach")
				settings := install.MakeClusterSettings(install.SecureOption(true))
				c.Start(ctx, t.L(), option.DefaultStartOpts(), settings)
				db := c.Conn(ctx, t.L(), 1)
				defer db.Close()
				_, err := db.QueryContext(ctx, `SELECT 1`)
				require.NoError(t, err)
			},
		})
	}
	r.Add(registry.TestSpec{
		Name:    "smoketest/secure/multitenant",
		Owner:   registry.OwnerMultiTenant,
		Cluster: r.MakeClusterSpec(2),
		Run:     multitenantSmokeTest,
	})
}

// multitenantSmokeTest verifies that a secure sql pod can connect to kv server
// and that tenant is is properly transmitted via cert.
func multitenantSmokeTest(ctx context.Context, t test.Test, c cluster.Cluster) {
	c.Put(ctx, t.Cockroach(), "./cockroach")
	settings := install.MakeClusterSettings(install.SecureOption(true))
	c.Start(ctx, t.L(), option.DefaultStartOpts(), settings, c.Node(1))

	// make sure connections to kvserver work
	db := c.Conn(ctx, t.L(), 1)
	defer db.Close()
	_, err := db.QueryContext(ctx, `SELECT 1`)
	require.NoError(t, err)

	tenID := 11
	ten := createTenantNode(ctx, t, c, c.Node(1), tenID, 2, 8011, 9011)
	runner := sqlutils.MakeSQLRunner(c.Conn(ctx, t.L(), 1))
	runner.Exec(t, `SELECT crdb_internal.create_tenant($1)`, tenID)
	ten.start(ctx, t, c, "./cockroach")

	// this doesn't work yet, roachprod knows nothing about tenants
	// db = c.Conn(ctx, t.L(), 2)
	// defer db.Close()

	tdb, err := gosql.Open("postgres", ten.pgURL)
	require.NoError(t, err)
	_, err = tdb.QueryContext(ctx, `SELECT 1`)
	require.NoError(t, err)

	// init kv and check new database was done right
	cmd := fmt.Sprintf("./cockroach workload init kv '%s'", ten.secureURL())
	err = c.RunE(ctx, c.Node(2), cmd)
	require.NoError(t, err)

	sqlutils.MakeSQLRunner(db).CheckQueryResultsRetry(t, fmt.Sprintf(`
    SELECT count(*) > 0
    FROM crdb_internal.ranges
    WHERE start_pretty LIKE '/Tenant/%d/%%';
`, tenID), [][]string{{"true"}})
}
