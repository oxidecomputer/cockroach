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

import React from "react";
import { ColumnDescriptor, SortedTable } from "src/sortedtable";
import { Tooltip } from "@cockroachlabs/ui-components";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import {
  Duration,
  formatNumberForDisplay,
  longToInt,
  TimestampToMoment,
  RenderCount,
  DATE_FORMAT_24_UTC,
  Count,
} from "../../util";

export type PlanHashStats = cockroach.server.serverpb.StatementDetailsResponse.ICollectedStatementGroupedByPlanHash;
export class PlansSortedTable extends SortedTable<PlanHashStats> {}

const planDetailsColumnLabels = {
  lastExecTime: "Last Execution Time",
  avgExecTime: "Average Execution Time",
  execCount: "Execution Count",
  avgRowsRead: "Average Rows Read",
  fullScan: "Full Scan",
  distSQL: "Distributed",
  planGist: "Plan Gist",
  vectorized: "Vectorized",
};
export type PlanDetailsTableColumnKeys = keyof typeof planDetailsColumnLabels;

type PlanDetailsTableTitleType = {
  [key in PlanDetailsTableColumnKeys]: () => JSX.Element;
};

export const planDetailsTableTitles: PlanDetailsTableTitleType = {
  planGist: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"The Gist of the Explain Plan."}
      >
        {planDetailsColumnLabels.planGist}
      </Tooltip>
    );
  },
  lastExecTime: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"The last time this Plan was executed."}
      >
        {planDetailsColumnLabels.lastExecTime}
      </Tooltip>
    );
  },
  avgExecTime: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"The average execution time for this Plan."}
      >
        {planDetailsColumnLabels.avgExecTime}
      </Tooltip>
    );
  },
  execCount: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"The execution count for this Plan."}
      >
        {planDetailsColumnLabels.execCount}
      </Tooltip>
    );
  },
  avgRowsRead: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"The average of rows read by this Plan."}
      >
        {planDetailsColumnLabels.avgRowsRead}
      </Tooltip>
    );
  },
  fullScan: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"If the Plan executed a Full Scan."}
      >
        {planDetailsColumnLabels.fullScan}
      </Tooltip>
    );
  },
  distSQL: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"If the Plan was distributed."}
      >
        {planDetailsColumnLabels.distSQL}
      </Tooltip>
    );
  },
  vectorized: () => {
    return (
      <Tooltip
        style="tableTitle"
        placement="bottom"
        content={"If the Plan was vectorized."}
      >
        {planDetailsColumnLabels.vectorized}
      </Tooltip>
    );
  },
};

export function makeExplainPlanColumns(
  handleDetails: (plan: PlanHashStats) => void,
): ColumnDescriptor<PlanHashStats>[] {
  const duration = (v: number) => Duration(v * 1e9);
  const count = (v: number) => v.toFixed(1);
  return [
    {
      name: "planGist",
      title: planDetailsTableTitles.planGist(),
      cell: (item: PlanHashStats) => (
        <Tooltip placement="bottom" content={item.stats.plan_gists[0]}>
          <a onClick={() => handleDetails(item)}>
            {item.stats.plan_gists[0].length > 25
              ? item.stats.plan_gists[0].slice(0, 22).concat("...")
              : item.stats.plan_gists[0]}
          </a>
        </Tooltip>
      ),
      sort: (item: PlanHashStats) => item.stats.plan_gists[0],
      alwaysShow: true,
    },
    {
      name: "lastExecTime",
      title: planDetailsTableTitles.lastExecTime(),
      cell: (item: PlanHashStats) =>
        TimestampToMoment(item.stats.last_exec_timestamp).format(
          DATE_FORMAT_24_UTC,
        ),
      sort: (item: PlanHashStats) =>
        TimestampToMoment(item.stats.last_exec_timestamp).unix(),
    },
    {
      name: "avgExecTime",
      title: planDetailsTableTitles.avgExecTime(),
      cell: (item: PlanHashStats) =>
        formatNumberForDisplay(item.stats.run_lat.mean, duration),
      sort: (item: PlanHashStats) => item.stats.run_lat.mean,
    },
    {
      name: "execCount",
      title: planDetailsTableTitles.execCount(),
      cell: (item: PlanHashStats) => Count(longToInt(item.stats.count)),
      sort: (item: PlanHashStats) => longToInt(item.stats.count),
    },
    {
      name: "avgRowsRead",
      title: planDetailsTableTitles.avgRowsRead(),
      cell: (item: PlanHashStats) =>
        formatNumberForDisplay(item.stats.rows_read.mean, count),
      sort: (item: PlanHashStats) =>
        formatNumberForDisplay(item.stats.rows_read.mean, count),
    },
    {
      name: "fullScan",
      title: planDetailsTableTitles.fullScan(),
      cell: (item: PlanHashStats) =>
        RenderCount(item.metadata.full_scan_count, item.metadata.total_count),
      sort: (item: PlanHashStats) =>
        RenderCount(item.metadata.full_scan_count, item.metadata.total_count),
    },
    {
      name: "distSQL",
      title: planDetailsTableTitles.distSQL(),
      cell: (item: PlanHashStats) =>
        RenderCount(item.metadata.dist_sql_count, item.metadata.total_count),
      sort: (item: PlanHashStats) =>
        RenderCount(item.metadata.dist_sql_count, item.metadata.total_count),
    },
    {
      name: "vectorized",
      title: planDetailsTableTitles.vectorized(),
      cell: (item: PlanHashStats) =>
        RenderCount(item.metadata.vec_count, item.metadata.total_count),
      sort: (item: PlanHashStats) =>
        RenderCount(item.metadata.vec_count, item.metadata.total_count),
    },
  ];
}
