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

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { AlignedData } from "uplot";
import { longToInt, TimestampToNumber } from "../util";

type statementStatisticsPerAggregatedTs = cockroach.server.serverpb.StatementDetailsResponse.ICollectedStatementGroupedByAggregatedTs;

export function generateExecuteAndPlanningTimeseries(
  stats: statementStatisticsPerAggregatedTs[],
): AlignedData {
  const ts: Array<number> = [];
  const execution: Array<number> = [];
  const planning: Array<number> = [];

  stats.forEach(function(stat: statementStatisticsPerAggregatedTs) {
    ts.push(TimestampToNumber(stat.aggregated_ts) * 1e3);
    execution.push(stat.stats.run_lat.mean * 1e9);
    planning.push(stat.stats.plan_lat.mean * 1e9);
  });

  return [ts, execution, planning];
}

export function generateRowsProcessedTimeseries(
  stats: statementStatisticsPerAggregatedTs[],
): AlignedData {
  const ts: Array<number> = [];
  const read: Array<number> = [];
  const written: Array<number> = [];

  stats.forEach(function(stat: statementStatisticsPerAggregatedTs) {
    ts.push(TimestampToNumber(stat.aggregated_ts) * 1e3);
    read.push(stat.stats.rows_read?.mean);
    written.push(stat.stats.rows_written?.mean);
  });

  return [ts, read, written];
}

export function generateExecRetriesTimeseries(
  stats: statementStatisticsPerAggregatedTs[],
): AlignedData {
  const ts: Array<number> = [];
  const retries: Array<number> = [];

  stats.forEach(function(stat: statementStatisticsPerAggregatedTs) {
    ts.push(TimestampToNumber(stat.aggregated_ts) * 1e3);

    const totalCountBarChart = longToInt(stat.stats.count);
    const firstAttemptsBarChart = longToInt(stat.stats.first_attempt_count);
    retries.push(totalCountBarChart - firstAttemptsBarChart);
  });

  return [ts, retries];
}

export function generateExecCountTimeseries(
  stats: statementStatisticsPerAggregatedTs[],
): AlignedData {
  const ts: Array<number> = [];
  const count: Array<number> = [];

  stats.forEach(function(stat: statementStatisticsPerAggregatedTs) {
    ts.push(TimestampToNumber(stat.aggregated_ts) * 1e3);
    count.push(longToInt(stat.stats.count));
  });

  return [ts, count];
}

export function generateContentionTimeseries(
  stats: statementStatisticsPerAggregatedTs[],
): AlignedData {
  const ts: Array<number> = [];
  const count: Array<number> = [];

  stats.forEach(function(stat: statementStatisticsPerAggregatedTs) {
    ts.push(TimestampToNumber(stat.aggregated_ts) * 1e3);
    count.push(stat.stats.exec_stats.contention_time.mean * 1e9);
  });

  return [ts, count];
}
