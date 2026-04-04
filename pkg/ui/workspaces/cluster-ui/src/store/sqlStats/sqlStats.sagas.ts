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

import { PayloadAction } from "@reduxjs/toolkit";
import { all, call, put, takeLatest, takeEvery } from "redux-saga/effects";
import {
  getCombinedStatements,
  StatementsRequest,
} from "src/api/statementsApi";
import { resetSQLStats } from "src/api/sqlStatsApi";
import { actions as localStorageActions } from "src/store/localStorage";
import {
  actions as sqlStatsActions,
  UpdateTimeScalePayload,
} from "./sqlStats.reducer";
import { actions as txnStatsActions } from "../transactionStats";
import { actions as sqlDetailsStatsActions } from "../statementDetails/statementDetails.reducer";

export function* refreshSQLStatsSaga(action: PayloadAction<StatementsRequest>) {
  yield put(sqlStatsActions.request(action.payload));
}

export function* requestSQLStatsSaga(
  action: PayloadAction<StatementsRequest>,
): any {
  try {
    const result = yield call(getCombinedStatements, action.payload);
    yield put(sqlStatsActions.received(result));
  } catch (e) {
    yield put(sqlStatsActions.failed(e));
  }
}

export function* updateSQLStatsTimeScaleSaga(
  action: PayloadAction<UpdateTimeScalePayload>,
) {
  const { ts } = action.payload;
  yield put(
    localStorageActions.updateTimeScale({
      value: ts,
    }),
  );
}

export function* resetSQLStatsSaga() {
  try {
    yield call(resetSQLStats);
    yield all([
      put(sqlDetailsStatsActions.invalidateAll()),
      put(sqlStatsActions.invalidated()),
      put(txnStatsActions.invalidated()),
    ]);
  } catch (e) {
    yield put(sqlStatsActions.failed(e));
  }
}

export function* sqlStatsSaga() {
  yield all([
    takeLatest(sqlStatsActions.refresh, refreshSQLStatsSaga),
    takeLatest(sqlStatsActions.request, requestSQLStatsSaga),
    takeLatest(sqlStatsActions.updateTimeScale, updateSQLStatsTimeScaleSaga),
    takeEvery(sqlStatsActions.reset, resetSQLStatsSaga),
  ]);
}
