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

package execinfra

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfrapb"
	"github.com/cockroachdb/cockroach/pkg/util/optional"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
	pbtypes "github.com/gogo/protobuf/types"
)

// ShouldCollectStats is a helper function used to determine if a processor
// should collect stats. The two requirements are that tracing must be enabled
// (to be able to output the stats somewhere), and that the flowCtx.CollectStats
// flag was set by the gateway node.
func ShouldCollectStats(ctx context.Context, flowCtx *FlowCtx) bool {
	return tracing.SpanFromContext(ctx) != nil && flowCtx.CollectStats
}

// GetCumulativeContentionTime is a helper function to calculate the cumulative
// contention time from the given recording or, if the recording is nil, from
// the tracing span from the context. All contention events found in the trace
// are included.
func GetCumulativeContentionTime(ctx context.Context, recording tracing.Recording) time.Duration {
	var cumulativeContentionTime time.Duration
	if recording == nil {
		recording = GetTraceData(ctx)
	}
	var ev roachpb.ContentionEvent
	for i := range recording {
		recording[i].Structured(func(any *pbtypes.Any, _ time.Time) {
			if !pbtypes.Is(any, &ev) {
				return
			}
			if err := pbtypes.UnmarshalAny(any, &ev); err != nil {
				return
			}
			cumulativeContentionTime += ev.Duration
		})
	}
	return cumulativeContentionTime
}

// ScanStats contains statistics on the internal MVCC operators used to satisfy
// a scan. See storage/engine.go for a more thorough discussion of the meaning
// of each stat.
// TODO(sql-observability): include other fields that are in roachpb.ScanStats,
// here and in execinfrapb.KVStats.
type ScanStats struct {
	// NumInterfaceSteps is the number of times the MVCC step function was called
	// to satisfy a scan.
	NumInterfaceSteps uint64
	// NumInternalSteps is the number of times that MVCC step was invoked
	// internally, including to step over internal, uncompacted Pebble versions.
	NumInternalSteps uint64
	// NumInterfaceSeeks is the number of times the MVCC seek function was called
	// to satisfy a scan.
	NumInterfaceSeeks uint64
	// NumInternalSeeks is the number of times that MVCC seek was invoked
	// internally, including to step over internal, uncompacted Pebble versions.
	NumInternalSeeks uint64
}

// PopulateKVMVCCStats adds data from the input ScanStats to the input KVStats.
func PopulateKVMVCCStats(kvStats *execinfrapb.KVStats, ss *ScanStats) {
	kvStats.NumInterfaceSteps = optional.MakeUint(ss.NumInterfaceSteps)
	kvStats.NumInternalSteps = optional.MakeUint(ss.NumInternalSteps)
	kvStats.NumInterfaceSeeks = optional.MakeUint(ss.NumInterfaceSeeks)
	kvStats.NumInternalSeeks = optional.MakeUint(ss.NumInternalSeeks)
}

// GetScanStats is a helper function to calculate scan stats from the given
// recording or, if the recording is nil, from the tracing span from the
// context.
func GetScanStats(ctx context.Context, recording tracing.Recording) (ss ScanStats) {
	if recording == nil {
		recording = GetTraceData(ctx)
	}
	var ev roachpb.ScanStats
	for i := range recording {
		recording[i].Structured(func(any *pbtypes.Any, _ time.Time) {
			if !pbtypes.Is(any, &ev) {
				return
			}
			if err := pbtypes.UnmarshalAny(any, &ev); err != nil {
				return
			}

			ss.NumInterfaceSteps += ev.NumInterfaceSteps
			ss.NumInternalSteps += ev.NumInternalSteps
			ss.NumInterfaceSeeks += ev.NumInterfaceSeeks
			ss.NumInternalSeeks += ev.NumInternalSeeks
		})
	}
	return ss
}
