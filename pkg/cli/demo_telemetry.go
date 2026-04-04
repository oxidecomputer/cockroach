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

package cli

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
)

// demoTelemetry corresponds to different sources of telemetry we are recording from cockroach demo.
type demoTelemetry int

const (
	_ demoTelemetry = iota
	// demo represents when cockroach demo is used at all.
	demo
	// nodes represents when cockroach demo is started with multiple nodes.
	nodes
	// demoLocality represents when cockroach demo is started with user defined localities.
	demoLocality
	// withLoad represents when cockroach demo is used with a background workload
	withLoad
	// geoPartitionedReplicas is used when cockroach demo is started with the geo-partitioned-replicas topology.
	geoPartitionedReplicas
)

var demoTelemetryMap = map[demoTelemetry]string{
	demo:                   "demo",
	nodes:                  "nodes",
	demoLocality:           "demo-locality",
	withLoad:               "withload",
	geoPartitionedReplicas: "geo-partitioned-replicas",
}

var demoTelemetryCounters map[demoTelemetry]telemetry.Counter

func init() {
	demoTelemetryCounters = make(map[demoTelemetry]telemetry.Counter)
	for ty, s := range demoTelemetryMap {
		demoTelemetryCounters[ty] = telemetry.GetCounterOnce(fmt.Sprintf("cli.demo.%s", s))
	}
}

func incrementDemoCounter(d demoTelemetry) {
	telemetry.Inc(demoTelemetryCounters[d])
}
