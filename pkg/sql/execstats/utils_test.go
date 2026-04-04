// Copyright 2020 The Cockroach Authors.
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

package execstats

import (
	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
)

// AddComponentStats modifies TraceAnalyzer internal state to add stats for the
// processor/stream/flow specified in stats.ComponentID and the given node ID.
func (a *TraceAnalyzer) AddComponentStats(stats *execinfrapb.ComponentStats) {
	a.FlowsMetadata.AddComponentStats(stats)
}

// AddComponentStats modifies FlowsMetadata to add stats for the
// processor/stream/flow specified in stats.ComponentID and the given node ID.
func (m *FlowsMetadata) AddComponentStats(stats *execinfrapb.ComponentStats) {
	switch stats.Component.Type {
	case execinfrapb.ComponentID_PROCESSOR:
		if m.processorStats == nil {
			m.processorStats = make(map[execinfrapb.ProcessorID]*execinfrapb.ComponentStats)
		}
		m.processorStats[execinfrapb.ProcessorID(stats.Component.ID)] = stats
	case execinfrapb.ComponentID_STREAM:
		streamStat := &streamStats{
			originSQLInstanceID: stats.Component.SQLInstanceID,
			stats:               stats,
		}
		if m.streamStats == nil {
			m.streamStats = make(map[execinfrapb.StreamID]*streamStats)
		}
		m.streamStats[execinfrapb.StreamID(stats.Component.ID)] = streamStat
	default:
		flowStat := &flowStats{}
		flowStat.stats = append(flowStat.stats, stats)
		if m.flowStats == nil {
			m.flowStats = make(map[base.SQLInstanceID]*flowStats)
		}
		m.flowStats[stats.Component.SQLInstanceID] = flowStat
	}
}
