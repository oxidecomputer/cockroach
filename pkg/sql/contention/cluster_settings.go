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

package contention

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/settings"
)

// TxnIDResolutionInterval is the cluster setting that controls how often the
// Transaction ID Resolution is performed.
var TxnIDResolutionInterval = settings.RegisterDurationSetting(
	settings.TenantWritable,
	"sql.contention.event_store.resolution_interval",
	"the interval at which transaction fingerprint ID resolution is "+
		"performed (set to 0 to disable)",
	time.Second*30,
)

// StoreCapacity is the cluster setting that controls the
// maximum size of the contention event store.
var StoreCapacity = settings.RegisterByteSizeSetting(
	settings.TenantWritable,
	"sql.contention.event_store.capacity",
	"the in-memory storage capacity per-node of contention event store",
	64*1024*1024, // 64 MB per node.
).WithPublic()

// DurationThreshold is the cluster setting for the threshold of
// contention durations. Only the contention events whose duration exceeds the
// threshold will be collected into crdb_internal.transaction_contention_events.
var DurationThreshold = settings.RegisterDurationSetting(
	settings.TenantWritable,
	"sql.contention.event_store.duration_threshold",
	"minimum contention duration to cause the contention events to be collected "+
		"into crdb_internal.transaction_contention_events",
	0,
).WithPublic()
