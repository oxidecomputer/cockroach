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

const addAvgSizeCol = `
ALTER TABLE system.table_statistics
ADD COLUMN IF NOT EXISTS "avgSize" INT8 NOT NULL DEFAULT (INT8 '0')
FAMILY "fam_0_tableID_statisticID_name_columnIDs_createdAt_rowCount_distinctCount_nullCount_histogram"
`

func alterSystemTableStatisticsAddAvgSize(
	ctx context.Context, cs clusterversion.ClusterVersion, d migration.TenantDeps, _ *jobs.Job,
) error {
	op := operation{
		name:           "add-table-statistics-avgSize-col",
		schemaList:     []string{"total_consumption"},
		query:          addAvgSizeCol,
		schemaExistsFn: hasColumn,
	}
	if err := migrateTable(ctx, cs, d, op, keys.TableStatisticsTableID, systemschema.TableStatisticsTable); err != nil {
		return err
	}
	return nil
}
