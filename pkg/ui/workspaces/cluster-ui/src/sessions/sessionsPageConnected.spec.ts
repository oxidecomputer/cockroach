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

import assert from "assert";
import fetchMock from "jest-fetch-mock";
import { applyMiddleware, createStore, Store } from "redux";
import createSagaMiddleware from "redux-saga";

import { rootReducer, sagas } from "src/store";
import {
  actions,
  ICancelQueryRequest,
  ICancelSessionRequest,
} from "src/store/terminateQuery";

class TestDriver {
  private readonly store: Store;

  constructor() {
    const sagaMiddleware = createSagaMiddleware();
    this.store = createStore(rootReducer, {}, applyMiddleware(sagaMiddleware));
    sagaMiddleware.run(sagas);
  }

  async cancelQuery(req: ICancelQueryRequest) {
    return this.store.dispatch(actions.terminateQuery(req));
  }

  async cancelSession(req: ICancelSessionRequest) {
    return this.store.dispatch(actions.terminateSession(req));
  }
}

describe("SessionsPage Connections", () => {
  beforeAll(fetchMock.enableMocks);
  afterEach(fetchMock.resetMocks);
  afterAll(fetchMock.disableMocks);

  describe("cancelQuery", () => {
    it("fires off an HTTP request", async () => {
      const driver = new TestDriver();
      assert.deepStrictEqual(fetchMock.mock.calls.length, 0);
      await driver.cancelQuery({ node_id: "1" });
      assert.deepStrictEqual(
        fetchMock.mock.calls[0][0],
        "/_status/cancel_query/1",
      );
    });
  });

  describe("cancelSession", () => {
    it("fires off an HTTP request", async () => {
      const driver = new TestDriver();
      assert.deepStrictEqual(fetchMock.mock.calls.length, 0);
      await driver.cancelSession({ node_id: "1" });
      assert.deepStrictEqual(
        fetchMock.mock.calls[0][0],
        "/_status/cancel_session/1",
      );
    });
  });
});
