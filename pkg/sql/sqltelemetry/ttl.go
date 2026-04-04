// Copyright 2022 The Cockroach Authors.
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

package sqltelemetry

import "github.com/cockroachdb/cockroach/pkg/server/telemetry"

var (
	// RowLevelTTLCreated is incremented when a row level TTL table is created.
	RowLevelTTLCreated = telemetry.GetCounterOnce("sql.row_level_ttl.created")

	// RowLevelTTLDropped is incremented when a row level TTL has been dropped
	// from a table.
	RowLevelTTLDropped = telemetry.GetCounterOnce("sql.row_level_ttl.dropped")

	// RowLevelTTLExecuted is incremented when a row level TTL job has executed.
	RowLevelTTLExecuted = telemetry.GetCounterOnce("sql.row_level_ttl.job_executed")
)
