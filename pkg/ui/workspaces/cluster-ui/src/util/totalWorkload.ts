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

import * as protos from "@cockroachlabs/crdb-protobuf-client";
import { AggregateStatistics } from "src/statementsTable";
import { longToInt } from "./fixLong";

type Statement = protos.cockroach.server.serverpb.StatementsResponse.ICollectedStatementStatistics;
type statementType = AggregateStatistics | Statement;
type statementsType = Array<statementType>;

/**
 * Function to calculate total workload of statements
 * Currently is recalculating every time is called, if that becomes an issue
 * on the future, consider use of cache and memoize the function
 * @param statements array of statements (AggregateStatistics or Statement)
 * @returns the total workload of all statements
 */
export function calculateTotalWorkload(statements: statementsType) {
  return statements.reduce((totalWorkload: number, stmt: statementType) => {
    return (totalWorkload +=
      longToInt(stmt.stats.count) * stmt.stats.service_lat.mean);
  }, 0);
}
