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

import { expectSaga } from "redux-saga-test-plan";
import { throwError } from "redux-saga-test-plan/providers";
import * as matchers from "redux-saga-test-plan/matchers";
import { getNodes } from "src/api/nodesApi";
import {
  receivedNodesSaga,
  requestNodesSaga,
  refreshNodesSaga,
} from "./nodes.sagas";
import { actions, reducer, NodesState } from "./nodes.reducer";
import { getNodesResponse } from "./nodes.fixtures";
import { accumulateMetrics } from "../../util";

describe("Nodes sagas", () => {
  const nodesResponse = getNodesResponse();
  const nodes = accumulateMetrics(nodesResponse.nodes);

  describe("refreshNodesSaga", () => {
    it("dispatches request nodes action", () => {
      expectSaga(refreshNodesSaga)
        .put(actions.request())
        .run();
    });
  });

  describe("requestNodesSaga", () => {
    it("successfully requests nodes list", () => {
      expectSaga(requestNodesSaga)
        .provide([[matchers.call.fn(getNodes), nodesResponse]])
        .put(actions.received(nodes))
        .withReducer(reducer)
        .hasFinalState<NodesState>({
          data: nodes,
          lastError: null,
          valid: true,
        })
        .run();
    });

    it("returns error on failed request", () => {
      const error = new Error("Failed request");
      expectSaga(requestNodesSaga)
        .provide([[matchers.call.fn(getNodes), throwError(error)]])
        .put(actions.failed(error))
        .withReducer(reducer)
        .hasFinalState<NodesState>({
          data: null,
          lastError: error,
          valid: false,
        })
        .run();
    });
  });

  describe("receivedNodesSaga", () => {
    it("sets valid status to false after specified period of time", () => {
      const timeout = 500;
      expectSaga(receivedNodesSaga, timeout)
        .delay(timeout)
        .put(actions.invalidated())
        .withReducer(reducer, {
          data: nodes,
          lastError: null,
          valid: true,
        })
        .hasFinalState<NodesState>({
          data: nodes,
          lastError: null,
          valid: false,
        })
        .run(1000);
    });
  });
});
