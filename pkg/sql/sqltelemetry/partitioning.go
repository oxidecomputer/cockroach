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

// PartitioningTelemetryType is an enum used to represent the different
// partitioning related operations that we are recording telemetry for.
type PartitioningTelemetryType int

const (
	_ PartitioningTelemetryType = iota
	// AlterAllPartitions represents an ALTER ALL PARTITIONS
	// statement (ALTER PARTITION OF INDEX t@*)
	AlterAllPartitions
	// PartitionConstrainedScan represents when the optimizer was
	// able to use partitioning to constrain a scan.
	PartitionConstrainedScan
)

var partitioningTelemetryMap = map[PartitioningTelemetryType]string{
	AlterAllPartitions:       "alter-all-partitions",
	PartitionConstrainedScan: "partition-constrained-scan",
}

func (p PartitioningTelemetryType) String() string {
	return partitioningTelemetryMap[p]
}

var partitioningTelemetryCounters map[PartitioningTelemetryType]telemetry.Counter

func init() {
	partitioningTelemetryCounters = make(map[PartitioningTelemetryType]telemetry.Counter)
	for ty, s := range partitioningTelemetryMap {
		partitioningTelemetryCounters[ty] = telemetry.GetCounterOnce(fmt.Sprintf("sql.partitioning.%s", s))
	}
}

// IncrementPartitioningCounter is used to increment the telemetry
// counter for a particular partitioning operation.
func IncrementPartitioningCounter(partitioningType PartitioningTelemetryType) {
	telemetry.Inc(partitioningTelemetryCounters[partitioningType])
}
