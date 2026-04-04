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

package sql_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// TestDestroyTenantSynchronous tests that destroying a tenant synchronously
// with crdb_internal.destroy_tenant(<ten_id>, true) works.
func TestDestroyTenantSynchronous(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	s, sqlDB, kvDB := serverutils.StartServer(t, base.TestServerArgs{})
	defer s.Stopper().Stop(ctx)

	tenantID := roachpb.MakeTenantID(10)
	codec := keys.MakeSQLCodec(tenantID)
	const tenantStateQuery = `
SELECT id, active FROM system.tenants WHERE id = 10
`
	checkKVsExistForTenant := func(t *testing.T, shouldExist bool) {
		rows, err := kvDB.Scan(
			ctx, codec.TenantPrefix(), codec.TenantPrefix().PrefixEnd(),
			1, /* maxRows */
		)
		require.NoError(t, err)
		require.Equal(t, shouldExist, len(rows) > 0)
	}

	tdb := sqlutils.MakeSQLRunner(sqlDB)

	// Create the tenant, make sure it has data and state.
	tdb.Exec(t, "SELECT crdb_internal.create_tenant(10)")
	tdb.CheckQueryResults(t, tenantStateQuery, [][]string{{"10", "true"}})
	checkKVsExistForTenant(t, true /* shouldExist */)

	// Destroy the tenant, make sure it does not have data and state.
	tdb.Exec(t, "SELECT crdb_internal.destroy_tenant(10, true)")
	tdb.CheckQueryResults(t, tenantStateQuery, [][]string{})
	checkKVsExistForTenant(t, false /* shouldExist */)
}
