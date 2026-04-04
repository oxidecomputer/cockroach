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

package sql

import (
	"sync/atomic"
	"time"

	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

// Default value used to designate the maximum frequency at which events
// are logged to the telemetry channel.
const defaultMaxEventFrequency = 8

var telemetryMaxEventFrequency = settings.RegisterIntSetting(
	settings.TenantWritable,
	"sql.telemetry.query_sampling.max_event_frequency",
	"the max event frequency at which we sample queries for telemetry, "+
		"note that this value shares a log-line limit of 10 logs per second on the "+
		"telemetry pipeline with all other telemetry events",
	defaultMaxEventFrequency,
	settings.NonNegativeInt,
)

// TelemetryLoggingMetrics keeps track of the last time at which an event
// was logged to the telemetry channel, and the number of skipped queries
// since the last logged event.
type TelemetryLoggingMetrics struct {
	mu struct {
		syncutil.RWMutex
		// The timestamp of the last emitted telemetry event.
		lastEmittedTime time.Time
	}
	Knobs *TelemetryLoggingTestingKnobs

	// skippedQueryCount is used to produce the count of non-sampled queries.
	skippedQueryCount uint64
}

// TelemetryLoggingTestingKnobs provides hooks and knobs for unit tests.
type TelemetryLoggingTestingKnobs struct {
	// getTimeNow allows tests to override the timeutil.Now() function used
	// when updating rolling query counts.
	getTimeNow func() time.Time
}

// ModuleTestingKnobs implements base.ModuleTestingKnobs interface.
func (*TelemetryLoggingTestingKnobs) ModuleTestingKnobs() {}

func (t *TelemetryLoggingMetrics) timeNow() time.Time {
	if t.Knobs != nil && t.Knobs.getTimeNow != nil {
		return t.Knobs.getTimeNow()
	}
	return timeutil.Now()
}

// maybeUpdateLastEmittedTime updates the lastEmittedTime if the amount of time
// elapsed between lastEmittedTime and newTime is greater than requiredSecondsElapsed.
func (t *TelemetryLoggingMetrics) maybeUpdateLastEmittedTime(
	newTime time.Time, requiredSecondsElapsed float64,
) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	lastEmittedTime := t.mu.lastEmittedTime

	if float64(newTime.Sub(lastEmittedTime))*1e-9 >= requiredSecondsElapsed {
		t.mu.lastEmittedTime = newTime
		return true
	}

	return false
}

func (t *TelemetryLoggingMetrics) resetSkippedQueryCount() (res uint64) {
	return atomic.SwapUint64(&t.skippedQueryCount, 0)
}

func (t *TelemetryLoggingMetrics) incSkippedQueryCount() {
	atomic.AddUint64(&t.skippedQueryCount, 1)
}
