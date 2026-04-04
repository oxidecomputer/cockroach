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
import { getNodes } from "src/api/nodesApi";
import { actions } from "./nodes.reducer";

import { CACHE_INVALIDATION_PERIOD, throttleWithReset } from "src/store/utils";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { rootActions } from "../reducers";

export function* refreshNodesSaga() {
  yield put(actions.request());
}

export function* requestNodesSaga() {
  try {
    const result: cockroach.server.serverpb.NodesResponse = yield call(
      getNodes,
    );
    yield put(actions.received(result.nodes));
  } catch (e) {
    yield put(actions.failed(e));
  }
}

export function* receivedNodesSaga(delayMs: number) {
  yield delay(delayMs);
  yield put(actions.invalidated());
}

export function* nodesSaga(
  cacheInvalidationPeriod: number = CACHE_INVALIDATION_PERIOD,
) {
  yield all([
    throttleWithReset(
      cacheInvalidationPeriod,
      actions.refresh,
      [actions.invalidated, rootActions.resetState],
      refreshNodesSaga,
    ),
    takeLatest(actions.request, requestNodesSaga),
    takeLatest(actions.received, receivedNodesSaga, cacheInvalidationPeriod),
  ]);
}
