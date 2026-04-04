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

import "github.com/cockroachdb/cockroach/pkg/util/metric"

// Metrics is a struct that include all metrics related to contention event
// store.
type Metrics struct {
	ResolverQueueSize *metric.Gauge
	ResolverRetries   *metric.Counter
	ResolverFailed    *metric.Counter
}

var _ metric.Struct = Metrics{}

// MetricStruct returns a new instance of Metrics.
func (Metrics) MetricStruct() {}

// NewMetrics returns a new instance of Metrics.
func NewMetrics() Metrics {
	return Metrics{
		ResolverQueueSize: metric.NewGauge(metric.Metadata{
			Name:        "sql.contention.resolver.queue_size",
			Help:        "Length of queued unresolved contention events",
			Measurement: "Queue length",
			Unit:        metric.Unit_COUNT,
		}),
		ResolverRetries: metric.NewCounter(metric.Metadata{
			Name:        "sql.contention.resolver.retries",
			Help:        "Number of times transaction id resolution has been retried",
			Measurement: "Retry count",
			Unit:        metric.Unit_COUNT,
		}),
		ResolverFailed: metric.NewCounter(metric.Metadata{
			Name:        "sql.contention.resolver.failed_resolutions",
			Help:        "Number of failed transaction ID resolution attempts",
			Measurement: "Failed transaction ID resolution count",
			Unit:        metric.Unit_COUNT,
		}),
	}
}
