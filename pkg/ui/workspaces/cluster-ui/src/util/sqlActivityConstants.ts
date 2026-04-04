// Copyright 2023 The Cockroach Authors.
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

import { duration } from "moment";
import { SqlStatsSortOptions, SqlStatsSortType } from "src/api/statementsApi";

export const limitOptions = [
  { value: 25, label: "25" },
  { value: 50, label: "50" },
  { value: 100, label: "100" },
  { value: 500, label: "500" },
];

export function getSortLabel(sort: SqlStatsSortType): string {
  switch (sort) {
    case SqlStatsSortOptions.SERVICE_LAT:
      return "Service Latency";
    case SqlStatsSortOptions.EXECUTION_COUNT:
      return "Execution Count";
    case SqlStatsSortOptions.CONTENTION_TIME:
      return "Contention Time";
    case SqlStatsSortOptions.PCT_RUNTIME:
      return "% Of All Run Time";
    default:
      return "";
  }
}

export const stmtRequestSortOptions = Object.values(SqlStatsSortOptions).map(
  sortVal => ({
    value: sortVal as SqlStatsSortType,
    label: getSortLabel(sortVal as SqlStatsSortType),
  }),
);

export const txnRequestSortOptions = stmtRequestSortOptions.filter(
  option => option.value !== SqlStatsSortOptions.PCT_RUNTIME,
);

export const STATS_LONG_LOADING_DURATION = duration(2, "s");
