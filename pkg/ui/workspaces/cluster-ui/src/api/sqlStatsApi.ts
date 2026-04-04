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

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { fetchData } from "src/api";

const RESET_SQL_STATS_PATH = "/_status/resetsqlstats";

export const resetSQLStats = (): Promise<cockroach.server.serverpb.ResetSQLStatsResponse> => {
  return fetchData(
    cockroach.server.serverpb.ResetSQLStatsResponse,
    RESET_SQL_STATS_PATH,
    cockroach.server.serverpb.ResetSQLStatsRequest,
    new cockroach.server.serverpb.ResetSQLStatsRequest({
      // reset_persisted_stats is set to true in order to clear both
      // in-memory stats as well as persisted stats.
      reset_persisted_stats: true,
    }),
  );
};
