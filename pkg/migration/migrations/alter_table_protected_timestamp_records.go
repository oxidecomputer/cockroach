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
	"github.com/cockroachdb/cockroach/pkg/migration"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/systemschema"
)

const addTargetCol = `
ALTER TABLE system.protected_ts_records
ADD COLUMN IF NOT EXISTS target BYTES FAMILY "primary"
`

func alterTableProtectedTimestampRecords(
	ctx context.Context, cs clusterversion.ClusterVersion, d migration.TenantDeps, _ *jobs.Job,
) error {
	op := operation{
		name:           "add-table-pts-records-target-col",
		schemaList:     []string{"target"},
		query:          addTargetCol,
		schemaExistsFn: hasColumn,
	}
	if err := migrateTable(ctx, cs, d, op,
		keys.ProtectedTimestampsRecordsTableID,
		systemschema.ProtectedTimestampsRecordsTable); err != nil {
		return err
	}
	return nil
}
