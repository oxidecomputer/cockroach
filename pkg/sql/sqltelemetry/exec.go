// Copyright 2019 The Cockroach Authors.
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

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
)

// DistSQLExecCounter is to be incremented whenever a query is distributed
// across multiple nodes.
var DistSQLExecCounter = telemetry.GetCounterOnce("sql.exec.query.is-distributed")

// VecExecCounter is to be incremented whenever a query runs with the vectorized
// execution engine.
var VecExecCounter = telemetry.GetCounterOnce("sql.exec.query.is-vectorized")

// VecModeCounter is to be incremented every time the vectorized execution mode
// is changed (including turned off).
func VecModeCounter(mode string) telemetry.Counter {
	return telemetry.GetCounter(fmt.Sprintf("sql.exec.vectorized-setting.%s", mode))
}

// CascadesLimitReached is to be incremented whenever the limit of foreign key
// cascade for a single query is exceeded.
var CascadesLimitReached = telemetry.GetCounterOnce("sql.exec.cascade-limit-reached")

// HashAggregationDiskSpillingDisabled is to be incremented whenever the disk
// spilling of the vectorized hash aggregator is disabled.
var HashAggregationDiskSpillingDisabled = telemetry.GetCounterOnce("sql.exec.hash-agg-spilling-disabled")

// DistSQLFlowsScheduled is to be incremented whenever a remote DistSQL flow is
// scheduled for running (regardless of whether it is being run right away or
// queued).
var DistSQLFlowsScheduled = telemetry.GetCounterOnce("sql.distsql.flows.scheduled")

// DistSQLFlowsQueued is to be incremented whenever a remote DistSQL flow is
// queued rather is run right away (because the node has reached
// 'sql.distsql.max_running_flows' limit).
var DistSQLFlowsQueued = telemetry.GetCounterOnce("sql.distsql.flows.queued")
