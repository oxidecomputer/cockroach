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

import { PayloadAction } from "@reduxjs/toolkit";
import { all, call, put, delay, takeLatest } from "redux-saga/effects";
import {
  ErrorWithKey,
  getStatementDetails,
  StatementDetailsRequest,
  StatementDetailsResponseWithKey,
} from "src/api/statementsApi";
import { actions as sqlDetailsStatsActions } from "./statementDetails.reducer";
import { CACHE_INVALIDATION_PERIOD } from "src/store/utils";
import { generateStmtDetailsToID } from "src/util/appStats";

export function* refreshSQLDetailsStatsSaga(
  action: PayloadAction<StatementDetailsRequest>,
) {
  yield put(sqlDetailsStatsActions.request(action?.payload));
}

export function* requestSQLDetailsStatsSaga(
  action: PayloadAction<StatementDetailsRequest>,
): any {
  const key = action?.payload
    ? generateStmtDetailsToID(
        action.payload.fingerprint_id,
        action.payload.app_names.toString(),
        action.payload.start,
        action.payload.end,
      )
    : "";
  try {
    const result = yield call(getStatementDetails, action?.payload);
    const resultWithKey: StatementDetailsResponseWithKey = {
      stmtResponse: result,
      key,
    };
    yield put(sqlDetailsStatsActions.received(resultWithKey));
  } catch (e) {
    const err: ErrorWithKey = {
      err: e,
      key,
    };
    yield put(sqlDetailsStatsActions.failed(err));
  }
}

export function receivedSQLDetailsStatsSagaFactory(delayMs: number) {
  return function* receivedSQLDetailsStatsSaga(
    action: PayloadAction<StatementDetailsResponseWithKey>,
  ) {
    yield delay(delayMs);
    yield put(
      sqlDetailsStatsActions.invalidated({
        key: action?.payload.key,
      }),
    );
  };
}

export function* sqlDetailsStatsSaga(
  cacheInvalidationPeriod: number = CACHE_INVALIDATION_PERIOD,
) {
  yield all([
    takeLatest(sqlDetailsStatsActions.refresh, refreshSQLDetailsStatsSaga),
    takeLatest(sqlDetailsStatsActions.request, requestSQLDetailsStatsSaga),
    takeLatest(
      sqlDetailsStatsActions.received,
      receivedSQLDetailsStatsSagaFactory(cacheInvalidationPeriod),
    ),
  ]);
}
