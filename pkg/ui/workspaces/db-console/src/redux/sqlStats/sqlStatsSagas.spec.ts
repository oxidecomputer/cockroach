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

import { expectSaga } from "redux-saga-test-plan";
import { call } from "redux-saga-test-plan/matchers";

import { resetSQLStatsFailedAction } from "./sqlStatsActions";
import { resetSQLStatsSaga } from "./sqlStatsSagas";
import { resetSQLStats } from "src/util/api";
import {
  apiReducersReducer,
  invalidateStatements,
  invalidateAllStatementDetails,
} from "src/redux/apiReducers";
import { throwError } from "redux-saga-test-plan/providers";

import { cockroach } from "src/js/protos";

describe("SQL Stats sagas", () => {
  describe("resetSQLStatsSaga", () => {
    const resetSQLStatsResponse = new cockroach.server.serverpb.ResetSQLStatsResponse();

    it("successfully resets SQL stats", () => {
      // TODO(azhng): validate refreshStatement() actions once we can figure out
      //  how to get ThunkAction to work with sagas.
      return expectSaga(resetSQLStatsSaga)
        .withReducer(apiReducersReducer)
        .provide([[call.fn(resetSQLStats), resetSQLStatsResponse]])
        .put(invalidateStatements())
        .put(invalidateAllStatementDetails())
        .run();
    });

    it("returns error on failed reset", () => {
      const err = new Error("failed to reset");
      return expectSaga(resetSQLStatsSaga)
        .provide([[call.fn(resetSQLStats), throwError(err)]])
        .put(resetSQLStatsFailedAction())
        .run();
    });
  });
});
