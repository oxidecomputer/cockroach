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

package txnidcache

import "github.com/cockroachdb/cockroach/pkg/util/metric"

// Metrics is a struct that include all metrics related to txn id cache.
type Metrics struct {
	CacheMissCounter *metric.Counter
	CacheReadCounter *metric.Counter
}

var _ metric.Struct = Metrics{}

// MetricStruct implements the metric.Struct interface.
func (Metrics) MetricStruct() {}

// NewMetrics returns a new instance of Metrics.
func NewMetrics() Metrics {
	return Metrics{
		CacheMissCounter: metric.NewCounter(metric.Metadata{
			Name:        "sql.contention.txn_id_cache.miss",
			Help:        "Number of cache misses",
			Measurement: "Cache miss",
			Unit:        metric.Unit_COUNT,
		}),
		CacheReadCounter: metric.NewCounter(metric.Metadata{
			Name:        "sql.contention.txn_id_cache.read",
			Help:        "Number of cache read",
			Measurement: "Cache read",
			Unit:        metric.Unit_COUNT,
		}),
	}
}
