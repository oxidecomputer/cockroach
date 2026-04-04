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

package migrations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/desctestutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/testcluster"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// TestEnsureNoDrainingNames tests if comments on indexes all have indexes,
// that exist.
func TestEnsureIndexesExistForComments(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	clusterArgs := base.TestClusterArgs{
		ServerArgs: base.TestServerArgs{
			Knobs: base.TestingKnobs{
				Server: &server.TestingKnobs{
					DisableAutomaticVersionUpgrade: make(chan struct{}),
					BinaryVersionOverride: clusterversion.ByKey(
						clusterversion.DeleteCommentsWithDroppedIndexes - 1),
				},
			},
		},
	}

	ctx := context.Background()
	tc := testcluster.StartTestCluster(t, 1, clusterArgs)
	s := tc.Server(0)
	defer tc.Stopper().Stop(ctx)
	sqlDB := tc.ServerConn(0)
	tdb := sqlutils.MakeSQLRunner(sqlDB)

	// Create a table, index and insert a valid comment and an artificial comment,
	// that belongs to an invalid index ID.
	tdb.Exec(t,
		"CREATE TABLE t1(name int);",
	)
	tdb.Exec(t,
		"CREATE INDEX blah ON t1(name);",
	)
	tdb.Exec(t,
		"COMMENT ON INDEX blah IS 'Valid comment';",
	)
	desc := desctestutils.TestingGetTableDescriptor(s.DB(), keys.SystemSQLCodec, "defaultdb", "public", "t1")
	tdb.Exec(t,
		"INSERT INTO system.comments VALUES($1, $2, 999, 'Invalid comment')",
		keys.IndexCommentType,
		desc.GetID(),
	)
	// Validate the state of the comments table we should have two entries.
	commentsTableStr := tdb.QueryStr(
		t,
		"SELECT * FROM system.comments ORDER BY sub_id ASC",
	)
	indexCommentTypeStr := fmt.Sprintf("%d", keys.IndexCommentType)
	descIDStr := fmt.Sprintf("%d", desc.GetID())
	require.Equal(t,
		commentsTableStr,
		[][]string{
			{indexCommentTypeStr, descIDStr, "2", "Valid comment"},
			{indexCommentTypeStr, descIDStr, "999", "Invalid comment"},
		},
	)

	// Migrate to the new cluster version.
	tdb.Exec(t, `SET CLUSTER SETTING version = $1`,
		clusterversion.ByKey(clusterversion.DeleteCommentsWithDroppedIndexes).String())

	tdb.CheckQueryResultsRetry(t, "SHOW CLUSTER SETTING version",
		[][]string{{clusterversion.ByKey(clusterversion.DeleteCommentsWithDroppedIndexes).String()}})

	// Validate the state of the comments table should only have the valid index
	// id.
	commentsTableStr = tdb.QueryStr(
		t,
		"SELECT * FROM system.comments ORDER BY sub_id ASC",
	)
	require.Equal(t,
		commentsTableStr,
		[][]string{
			{indexCommentTypeStr, descIDStr, "2", "Valid comment"},
		},
	)
}
