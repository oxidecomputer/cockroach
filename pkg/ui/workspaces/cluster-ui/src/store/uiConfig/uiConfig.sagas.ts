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

import { all, call, delay, put, takeLatest } from "redux-saga/effects";
import { actions } from "./uiConfig.reducer";
import { getUserSQLRoles } from "../../api/userApi";
import { CACHE_INVALIDATION_PERIOD, throttleWithReset } from "../utils";
import { rootActions } from "../reducers";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";

export function* refreshUserSQLRolesSaga(): any {
  yield put(actions.requestUserSQLRoles());
}

export function* requestUserSQLRolesSaga(): any {
  try {
    const result: cockroach.server.serverpb.UserSQLRolesResponse = yield call(
      getUserSQLRoles,
    );
    yield put(actions.receivedUserSQLRoles(result.roles));
  } catch (e) {
    console.warn(e.message);
  }
}

export function* receivedUserSQLRolesSaga(delayMs: number): any {
  yield delay(delayMs);
  yield put(actions.invalidatedUserSQLRoles());
}

export function* uiConfigSaga(
  cacheInvalidationPeriod: number = CACHE_INVALIDATION_PERIOD,
): any {
  yield all([
    throttleWithReset(
      cacheInvalidationPeriod,
      actions.refreshUserSQLRoles,
      [actions.invalidatedUserSQLRoles, rootActions.resetState],
      refreshUserSQLRolesSaga,
    ),
    takeLatest(actions.requestUserSQLRoles, requestUserSQLRolesSaga),
    takeLatest(
      actions.receivedUserSQLRoles,
      receivedUserSQLRolesSaga,
      cacheInvalidationPeriod,
    ),
  ]);
}
