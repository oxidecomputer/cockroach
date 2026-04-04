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

import { PayloadAction } from "@reduxjs/toolkit";
import { all, call, put, takeLatest } from "redux-saga/effects";
import {
  getFlushedTxnStatsApi,
  StatementsRequest,
} from "src/api/statementsApi";
import { actions as txnStatsActions } from "./txnStats.reducer";

export function* refreshTxnStatsSaga(
  action: PayloadAction<StatementsRequest>,
): any {
  yield put(txnStatsActions.request(action.payload));
}

export function* requestTxnStatsSaga(
  action: PayloadAction<StatementsRequest>,
): any {
  try {
    const result = yield call(getFlushedTxnStatsApi, action.payload);
    yield put(txnStatsActions.received(result));
  } catch (e) {
    yield put(txnStatsActions.failed(e));
  }
}

export function* txnStatsSaga(): any {
  yield all([
    takeLatest(txnStatsActions.refresh, refreshTxnStatsSaga),
    takeLatest(txnStatsActions.request, requestTxnStatsSaga),
  ]);
}
