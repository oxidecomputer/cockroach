// Copyright 2020 The Cockroach Authors.
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
import { all, call, put, takeEvery } from "redux-saga/effects";

import { terminateQuery, terminateSession } from "src/api/terminateQueryApi";
import { actions as sessionsActions } from "src/store/sessions";
import { actions as terminateQueryActions } from "./terminateQuery.reducer";

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";

const CancelSessionRequest = cockroach.server.serverpb.CancelSessionRequest;
const CancelQueryRequest = cockroach.server.serverpb.CancelQueryRequest;
export type ICancelSessionRequest = cockroach.server.serverpb.ICancelSessionRequest;
export type ICancelQueryRequest = cockroach.server.serverpb.ICancelQueryRequest;

export function* terminateSessionSaga(
  action: PayloadAction<ICancelSessionRequest>,
) {
  const terminateSessionRequest = new CancelSessionRequest(action.payload);
  try {
    yield call(terminateSession, terminateSessionRequest);
    yield put(terminateQueryActions.terminateSessionCompleted());
    yield put(sessionsActions.invalidated());
    yield put(sessionsActions.refresh());
  } catch (e) {
    yield put(terminateQueryActions.terminateSessionFailed(e));
  }
}

export function* terminateQuerySaga(
  action: PayloadAction<ICancelQueryRequest>,
) {
  const terminateQueryRequest = new CancelQueryRequest(action.payload);
  try {
    yield call(terminateQuery, terminateQueryRequest);
    yield put(terminateQueryActions.terminateQueryCompleted());
    yield put(sessionsActions.invalidated());
    yield put(sessionsActions.refresh());
  } catch (e) {
    yield put(terminateQueryActions.terminateQueryFailed(e));
  }
}

export function* terminateSaga() {
  yield all([
    takeEvery(terminateQueryActions.terminateSession, terminateSessionSaga),
    takeEvery(terminateQueryActions.terminateQuery, terminateQuerySaga),
  ]);
}
