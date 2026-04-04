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

package tenantrate

import (
	"github.com/cockroachdb/cockroach/pkg/multitenant"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
	"github.com/cockroachdb/cockroach/pkg/util/metric/aggmetric"
)

// Metrics is a metric.Struct for the LimiterFactory.
type Metrics struct {
	Tenants               *metric.Gauge
	CurrentBlocked        *aggmetric.AggGauge
	ReadBatchesAdmitted   *aggmetric.AggCounter
	WriteBatchesAdmitted  *aggmetric.AggCounter
	ReadRequestsAdmitted  *aggmetric.AggCounter
	WriteRequestsAdmitted *aggmetric.AggCounter
	ReadBytesAdmitted     *aggmetric.AggCounter
	WriteBytesAdmitted    *aggmetric.AggCounter
}

var _ metric.Struct = (*Metrics)(nil)

var (
	metaTenants = metric.Metadata{
		Name:        "kv.tenant_rate_limit.num_tenants",
		Help:        "Number of tenants currently being tracked",
		Measurement: "Tenants",
		Unit:        metric.Unit_COUNT,
	}
	metaCurrentBlocked = metric.Metadata{
		Name:        "kv.tenant_rate_limit.current_blocked",
		Help:        "Number of requests currently blocked by the rate limiter",
		Measurement: "Requests",
		Unit:        metric.Unit_COUNT,
	}
	metaReadBatchesAdmitted = metric.Metadata{
		Name:        "kv.tenant_rate_limit.read_batches_admitted",
		Help:        "Number of read batches admitted by the rate limiter",
		Measurement: "Requests",
		Unit:        metric.Unit_COUNT,
	}
	metaWriteBatchesAdmitted = metric.Metadata{
		Name:        "kv.tenant_rate_limit.write_batches_admitted",
		Help:        "Number of write batches admitted by the rate limiter",
		Measurement: "Requests",
		Unit:        metric.Unit_COUNT,
	}
	metaReadRequestsAdmitted = metric.Metadata{
		Name:        "kv.tenant_rate_limit.read_requests_admitted",
		Help:        "Number of read requests admitted by the rate limiter",
		Measurement: "Requests",
		Unit:        metric.Unit_COUNT,
	}
	metaWriteRequestsAdmitted = metric.Metadata{
		Name:        "kv.tenant_rate_limit.write_requests_admitted",
		Help:        "Number of write requests admitted by the rate limiter",
		Measurement: "Requests",
		Unit:        metric.Unit_COUNT,
	}
	metaReadBytesAdmitted = metric.Metadata{
		Name:        "kv.tenant_rate_limit.read_bytes_admitted",
		Help:        "Number of read bytes admitted by the rate limiter",
		Measurement: "Bytes",
		Unit:        metric.Unit_BYTES,
	}
	metaWriteBytesAdmitted = metric.Metadata{
		Name:        "kv.tenant_rate_limit.write_bytes_admitted",
		Help:        "Number of write bytes admitted by the rate limiter",
		Measurement: "Bytes",
		Unit:        metric.Unit_BYTES,
	}
)

func makeMetrics() Metrics {
	b := aggmetric.MakeBuilder(multitenant.TenantIDLabel)
	return Metrics{
		Tenants:               metric.NewGauge(metaTenants),
		CurrentBlocked:        b.Gauge(metaCurrentBlocked),
		ReadBatchesAdmitted:   b.Counter(metaReadBatchesAdmitted),
		WriteBatchesAdmitted:  b.Counter(metaWriteBatchesAdmitted),
		ReadRequestsAdmitted:  b.Counter(metaReadRequestsAdmitted),
		WriteRequestsAdmitted: b.Counter(metaWriteRequestsAdmitted),
		ReadBytesAdmitted:     b.Counter(metaReadBytesAdmitted),
		WriteBytesAdmitted:    b.Counter(metaWriteBytesAdmitted),
	}
}

// MetricStruct indicates that Metrics is a metric.Struct
func (m *Metrics) MetricStruct() {}

// tenantMetrics represent metrics for an individual tenant.
type tenantMetrics struct {
	currentBlocked        *aggmetric.Gauge
	readBatchesAdmitted   *aggmetric.Counter
	writeBatchesAdmitted  *aggmetric.Counter
	readRequestsAdmitted  *aggmetric.Counter
	writeRequestsAdmitted *aggmetric.Counter
	readBytesAdmitted     *aggmetric.Counter
	writeBytesAdmitted    *aggmetric.Counter
}

func (m *Metrics) tenantMetrics(tenantID roachpb.TenantID) tenantMetrics {
	tid := tenantID.String()
	return tenantMetrics{
		currentBlocked:        m.CurrentBlocked.AddChild(tid),
		readBatchesAdmitted:   m.ReadBatchesAdmitted.AddChild(tid),
		writeBatchesAdmitted:  m.WriteBatchesAdmitted.AddChild(tid),
		readRequestsAdmitted:  m.ReadRequestsAdmitted.AddChild(tid),
		writeRequestsAdmitted: m.WriteRequestsAdmitted.AddChild(tid),
		readBytesAdmitted:     m.ReadBytesAdmitted.AddChild(tid),
		writeBytesAdmitted:    m.WriteBytesAdmitted.AddChild(tid),
	}
}

func (tm *tenantMetrics) destroy() {
	tm.currentBlocked.Destroy()
	tm.readBatchesAdmitted.Destroy()
	tm.writeBatchesAdmitted.Destroy()
	tm.readRequestsAdmitted.Destroy()
	tm.writeRequestsAdmitted.Destroy()
	tm.readBytesAdmitted.Destroy()
	tm.writeBytesAdmitted.Destroy()
}
