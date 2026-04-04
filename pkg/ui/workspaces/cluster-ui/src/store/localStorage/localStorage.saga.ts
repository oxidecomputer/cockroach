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

import { AnyAction } from "redux";
import { all, call, takeEvery, takeLatest, put } from "redux-saga/effects";
import {
  actions,
  LocalStorageKeys,
  TypedPayload,
} from "./localStorage.reducer";
import { actions as sqlStatsActions } from "src/store/sqlStats";
import { PayloadAction } from "@reduxjs/toolkit";
import { TimeScale } from "src/timeScaleDropdown";

export function* updateLocalStorageItemSaga(action: AnyAction) {
  const { key, value } = action.payload;
  yield call(
    { context: localStorage, fn: localStorage.setItem },
    key,
    JSON.stringify(value),
  );
}

export function* updateTimeScale(
  action: PayloadAction<TypedPayload<TimeScale>>,
) {
  yield put(sqlStatsActions.invalidated());
  yield call(
    { context: localStorage, fn: localStorage.setItem },
    LocalStorageKeys.GLOBAL_TIME_SCALE,
    JSON.stringify(action.payload?.value),
  );
}

export function* localStorageSaga() {
  yield all([
    takeEvery(actions.update, updateLocalStorageItemSaga),
    takeLatest(actions.updateTimeScale, updateTimeScale),
  ]);
}
