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

package delegate

import (
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
)

func (d *delegator) delegateShowFullTableScans() (tree.Statement, error) {
	sqltelemetry.IncrementShowCounter(sqltelemetry.FullTableScans)
	const query = `
  SELECT 
    key AS query, count, rows_read_avg, bytes_read_avg, service_lat_avg, contention_time_avg, max_mem_usage_avg, network_bytes_avg, max_retries
  FROM crdb_internal.node_statement_statistics WHERE full_scan = TRUE ORDER BY count DESC`
	return parse(query)
}
