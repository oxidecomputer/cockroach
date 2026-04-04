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

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/migration"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
)

// ensureCommentsHaveNonDroppedIndexes cleans up any comments associated with
// indexes that no longer exist.
func ensureCommentsHaveNonDroppedIndexes(
	ctx context.Context, _ clusterversion.ClusterVersion, d migration.TenantDeps, _ *jobs.Job,
) error {
	return d.DB.Txn(ctx, func(ctx context.Context, txn *kv.Txn) error {
		// Delete the rows that don't belong to any indexes.
		_, err := d.InternalExecutor.QueryBufferedEx(
			ctx,
			"select-comments-with-missing-indexes",
			txn,
			sessiondata.InternalExecutorOverride{User: security.RootUserName()},
			`DELETE FROM system.comments
      WHERE type = $1
            AND (object_id, sub_id)
				NOT IN (
						SELECT (descriptor_id, index_id)
						  FROM crdb_internal.table_indexes
					);`,
			keys.IndexCommentType,
		)
		if err != nil {
			return err
		}
		return txn.Commit(ctx)
	})
}
