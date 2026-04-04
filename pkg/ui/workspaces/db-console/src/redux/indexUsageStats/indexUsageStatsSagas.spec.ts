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
import { call, select } from "redux-saga-test-plan/matchers";

import {
  resetIndexUsageStatsFailedAction,
  resetIndexUsageStatsCompleteAction,
  resetIndexUsageStatsAction,
} from "./indexUsageStatsActions";
import {
  resetIndexUsageStatsSaga,
  selectIndexStatsKeys,
} from "./indexUsageStatsSagas";
import { resetIndexUsageStats } from "src/util/api";
import { throwError } from "redux-saga-test-plan/providers";

import { cockroach } from "src/js/protos";

describe("Index Usage Stats sagas", () => {
  describe("resetIndexUsageStatsSaga", () => {
    const resetIndexUsageStatsResponse = new cockroach.server.serverpb.ResetIndexUsageStatsResponse();
    const action = resetIndexUsageStatsAction("database", "table");

    it("successfully resets index usage stats", () => {
      // TODO(lindseyjin): figure out how to test invalidate and refresh actions
      //  once we can figure out how to get ThunkAction to work with sagas.
      return expectSaga(resetIndexUsageStatsSaga, action)
        .provide([
          [call.fn(resetIndexUsageStats), resetIndexUsageStatsResponse],
          [select(selectIndexStatsKeys), ["database/table"]],
        ])
        .put(resetIndexUsageStatsCompleteAction())
        .dispatch(action)
        .run();
    });

    it("returns error on failed reset", () => {
      const err = new Error("failed to reset");
      return expectSaga(resetIndexUsageStatsSaga, action)
        .provide([[call.fn(resetIndexUsageStats), throwError(err)]])
        .put(resetIndexUsageStatsFailedAction())
        .dispatch(action)
        .run();
    });
  });
});
