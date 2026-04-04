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

import { cockroach } from "src/js/protos";
import { resetSQLStats } from "src/util/api";
import { all, call, put, takeEvery } from "redux-saga/effects";
import { RESET_SQL_STATS, resetSQLStatsFailedAction } from "./sqlStatsActions";
import {
  invalidateAllStatementDetails,
  invalidateStatements,
} from "src/redux/apiReducers";

import ResetSQLStatsRequest = cockroach.server.serverpb.ResetSQLStatsRequest;

export function* resetSQLStatsSaga() {
  const resetSQLStatsRequest = new ResetSQLStatsRequest({
    // reset_persisted_stats is set to true in order to clear both
    // in-memory stats as well as persisted stats.
    reset_persisted_stats: true,
  });
  try {
    yield call(resetSQLStats, resetSQLStatsRequest);
    yield all([
      put(invalidateStatements()),
      put(invalidateAllStatementDetails()),
    ]);
  } catch (e) {
    yield put(resetSQLStatsFailedAction());
  }
}

export function* sqlStatsSaga() {
  yield all([takeEvery(RESET_SQL_STATS, resetSQLStatsSaga)]);
}
