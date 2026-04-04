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

import React from "react";
import _ from "lodash";

import { LineGraph } from "src/views/cluster/components/linegraph";
import { Metric, Axis } from "src/views/shared/components/metricQuery";
import { AxisUnits } from "@cockroachlabs/cluster-ui";

import { GraphDashboardProps } from "./dashboardUtils";

export default function(props: GraphDashboardProps) {
  const { nodeSources } = props;

  const percentiles = ["p50", "p75", "p90", "p95", "p99"];

  return [
    <LineGraph title="Processing Rate" sources={nodeSources}>
      <Axis label="rows per second" units={AxisUnits.Count}>
        <Metric
          name="cr.node.jobs.row_level_ttl.rows_selected"
          title="rows selected"
          nonNegativeRate
        />
        <Metric
          name="cr.node.jobs.row_level_ttl.rows_deleted"
          title="rows deleted"
          nonNegativeRate
        />
      </Axis>
    </LineGraph>,
    <LineGraph title="Estimated Rows" sources={nodeSources}>
      <Axis label="row count" units={AxisUnits.Count}>
        <Metric
          name="cr.node.jobs.row_level_ttl.total_rows"
          title="approximate number of rows"
          nonNegativeRate
        />
        <Metric
          name="cr.node.jobs.row_level_ttl.total_expired_rows"
          title="approximate number of expired rows"
          nonNegativeRate
        />
      </Axis>
    </LineGraph>,
    <LineGraph
      title="Job Latency"
      sources={nodeSources}
      tooltip={`Latency of scanning and deleting within the job.`}
    >
      <Axis label="latency" units={AxisUnits.Duration}>
        {_.map(percentiles, p => (
          <>
            <Metric
              name={`cr.node.jobs.row_level_ttl.select_duration-${p}`}
              title={`scan latency (${p})`}
              downsampleMax
            />
            <Metric
              name={`cr.node.jobs.row_level_ttl.delete_duration-${p}`}
              title={`delete latency (${p})`}
              downsampleMax
            />
          </>
        ))}
      </Axis>
    </LineGraph>,
    <LineGraph
      title="Ranges in Progress"
      sources={nodeSources}
      tooltip={`Number of active ranges being processed by TTL.`}
    >
      <Axis label="range count" units={AxisUnits.Count}>
        <Metric
          name="cr.node.jobs.row_level_ttl.num_active_ranges"
          title="number of ranges being processed"
          nonNegativeRate
        />
      </Axis>
    </LineGraph>,
  ];
}
