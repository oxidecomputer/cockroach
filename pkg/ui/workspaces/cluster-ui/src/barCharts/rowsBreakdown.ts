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

import { scaleLinear } from "d3-scale";
import { stdDevLong } from "src/util";
import { formatTwoPlaces, normalizeClosedDomain } from "./utils";
import * as protos from "@cockroachlabs/crdb-protobuf-client";

type StatementStatistics = protos.cockroach.server.serverpb.StatementsResponse.ICollectedStatementStatistics;

export function rowsBreakdown(s: StatementStatistics) {
  const mean = s.stats.num_rows.mean;
  const sd = stdDevLong(s.stats.num_rows, s.stats.count);
  const domain = normalizeClosedDomain([0, mean + sd]);

  const scale = scaleLinear()
    .domain(domain)
    .range([0, 100]);

  return {
    rowsBarChart(meanRow?: boolean) {
      const spread = scale(sd + (sd > mean ? mean : sd));
      if (meanRow) {
        return formatTwoPlaces(mean);
      } else {
        return spread;
      }
    },
  };
}
