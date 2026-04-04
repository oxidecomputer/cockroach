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

package server

import (
	"context"
	"net/http"

	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/server/status"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

// Construct a handler responsible for serving the instant values of selected
// load metrics. These include user and system CPU time currently.
// TODO(knz): this should probably include memory usage too somehow.
func makeStatusLoadHandler(
	ctx context.Context, rsr *status.RuntimeStatSampler, metricSource metricMarshaler,
) func(http.ResponseWriter, *http.Request) {
	cpuUserNanos := metric.NewGauge(rsr.CPUUserNS.GetMetadata())
	cpuSysNanos := metric.NewGauge(rsr.CPUSysNS.GetMetadata())
	cpuNowNanos := metric.NewGauge(rsr.CPUNowNS.GetMetadata())
	registry := metric.NewRegistry()
	registry.AddMetric(cpuUserNanos)
	registry.AddMetric(cpuSysNanos)
	registry.AddMetric(cpuNowNanos)

	// Exporter for the CPU metrics that are provided only by the load handler.
	exporter := metric.MakePrometheusExporter()
	regScrape := func(pm *metric.PrometheusExporter) {
		pm.ScrapeRegistry(registry, true)
	}

	// Exporter for the selected metrics that also show in /_status/vars.
	exporter2 := metric.MakePrometheusExporterForSelectedMetrics(map[string]struct{}{
		sql.MetaQueryExecuted.Name:       {},
		pgwire.MetaConns.Name:            {},
		jobs.MetaRunningNonIdleJobs.Name: {},
	})

	return func(w http.ResponseWriter, r *http.Request) {
		userTimeMillis, sysTimeMillis, err := status.GetCPUTime(ctx)
		if err != nil {
			// Just log but don't return an error to match the _status/vars metrics handler.
			log.Ops.Errorf(ctx, "unable to get cpu usage: %v", err)
		}

		// The CPU metrics are updated on each call.
		// cpuTime.{User,Sys} are in milliseconds, convert to nanoseconds.
		utime := userTimeMillis * 1e6
		stime := sysTimeMillis * 1e6
		cpuUserNanos.Update(utime)
		cpuSysNanos.Update(stime)
		cpuNowNanos.Update(timeutil.Now().UnixNano())

		if err := exporter.ScrapeAndPrintAsText(w, regScrape); err != nil {
			log.Errorf(r.Context(), "%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := exporter2.ScrapeAndPrintAsText(w, metricSource.ScrapeIntoPrometheus); err != nil {
			log.Errorf(r.Context(), "%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
