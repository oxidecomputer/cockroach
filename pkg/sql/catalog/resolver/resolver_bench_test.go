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

package resolver_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/resolver"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondatapb"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/stretchr/testify/require"
)

// BenchmarkResolveExistingObject exercises resolver code to ensure that
func BenchmarkResolveExistingObject(b *testing.B) {
	defer leaktest.AfterTest(b)()
	defer log.Scope(b).Close(b)
	for _, tc := range []struct {
		setup      []string
		name       tree.UnresolvedName
		flags      tree.ObjectLookupFlags
		searchPath string
	}{
		{
			setup: []string{
				"CREATE TABLE foo ()",
			},
			name:  tree.MakeUnresolvedName("foo"),
			flags: tree.ObjectLookupFlagsWithRequired(),
		},
		{
			setup: []string{
				"CREATE SCHEMA sc",
				"CREATE TABLE sc.foo ()",
			},
			name:  tree.MakeUnresolvedName("sc", "foo"),
			flags: tree.ObjectLookupFlagsWithRequired(),
		},
		{
			setup: []string{
				"CREATE SCHEMA sc",
				"CREATE TABLE sc.foo ()",
			},
			name:       tree.MakeUnresolvedName("foo"),
			flags:      tree.ObjectLookupFlagsWithRequired(),
			searchPath: "public,$user,sc",
		},
	} {
		b.Run(strings.Join(tc.setup, ";")+tc.name.String(), func(b *testing.B) {
			ctx := context.Background()
			s, sqlDB, kvDB := serverutils.StartServer(b, base.TestServerArgs{})
			defer s.Stopper().Stop(ctx)
			tDB := sqlutils.MakeSQLRunner(sqlDB)
			for _, stmt := range tc.setup {
				tDB.Exec(b, stmt)
			}

			var sessionData sessiondatapb.SessionData
			{
				var sessionSerialized []byte
				tDB.QueryRow(b, "SELECT crdb_internal.serialize_session()").Scan(&sessionSerialized)
				require.NoError(b, protoutil.Unmarshal(sessionSerialized, &sessionData))
			}

			execCfg := s.ExecutorConfig().(sql.ExecutorConfig)
			txn := kvDB.NewTxn(ctx, "test")
			p, cleanup := sql.NewInternalPlanner("asdf", txn, security.RootUserName(), &sql.MemoryMetrics{}, &execCfg, sessionData)
			defer cleanup()

			// The internal planner overrides the database to "system", here we
			// change it back.
			{
				ec := p.(interface{ EvalContext() *tree.EvalContext }).EvalContext()
				require.NoError(b, ec.SessionAccessor.SetSessionVar(
					ctx, "database", "defaultdb", false,
				))

				if tc.searchPath != "" {
					require.NoError(b, ec.SessionAccessor.SetSessionVar(
						ctx, "search_path", tc.searchPath, false,
					))
				}
			}

			rs := p.(resolver.SchemaResolver)
			uon, err := tc.name.ToUnresolvedObjectName(0)
			require.NoError(b, err)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				desc, _, err := resolver.ResolveExistingObject(ctx, rs, uon, tc.flags)
				require.NoError(b, err)
				require.NotNil(b, desc)
			}
		})
	}
}
