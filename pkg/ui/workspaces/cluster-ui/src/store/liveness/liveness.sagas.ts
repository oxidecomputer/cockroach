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

import { all, call, put, delay, takeLatest } from "redux-saga/effects";
import { getLiveness } from "src/api/livenessApi";
import { actions } from "./liveness.reducer";

import { CACHE_INVALIDATION_PERIOD, throttleWithReset } from "src/store/utils";
import { rootActions } from "../reducers";

export function* refreshLivenessSaga() {
  yield put(actions.request());
}

export function* requestLivenessSaga(): any {
  try {
    const result = yield call(getLiveness);
    yield put(actions.received(result));
  } catch (e) {
    yield put(actions.failed(e));
  }
}

export function* receivedLivenessSaga(delayMs: number) {
  yield delay(delayMs);
  yield put(actions.invalidated());
}

export function* livenessSaga(
  cacheInvalidationPeriod: number = CACHE_INVALIDATION_PERIOD,
) {
  yield all([
    throttleWithReset(
      cacheInvalidationPeriod,
      actions.refresh,
      [actions.invalidated, rootActions.resetState],
      refreshLivenessSaga,
    ),
    takeLatest(actions.request, requestLivenessSaga),
    takeLatest(actions.received, receivedLivenessSaga, cacheInvalidationPeriod),
  ]);
}
