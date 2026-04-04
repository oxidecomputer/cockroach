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

import { all, call, delay, put, takeLatest } from "redux-saga/effects";

import { actions } from "./sessions.reducer";
import { getSessions } from "src/api/sessionsApi";
import { CACHE_INVALIDATION_PERIOD, throttleWithReset } from "../utils";
import { rootActions } from "../reducers";

export function* refreshSessionsSaga() {
  yield put(actions.request());
}

export function* requestSessionsSaga(): any {
  try {
    const result = yield call(getSessions);
    yield put(actions.received(result));
  } catch (e) {
    yield put(actions.failed(e));
  }
}

export function* receivedStatementsSaga(delayMs: number) {
  yield delay(delayMs);
  yield put(actions.invalidated());
}

export function* sessionsSaga(
  cacheInvalidationPeriod: number = CACHE_INVALIDATION_PERIOD,
) {
  yield all([
    throttleWithReset(
      cacheInvalidationPeriod,
      actions.refresh,
      [actions.invalidated, rootActions.resetState],
      refreshSessionsSaga,
    ),
    takeLatest(actions.request, requestSessionsSaga),
    takeLatest(
      actions.received,
      receivedStatementsSaga,
      cacheInvalidationPeriod,
    ),
  ]);
}
