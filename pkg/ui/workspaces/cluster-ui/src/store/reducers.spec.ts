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

import { assert } from "chai";
import { createStore } from "redux";
import { rootActions, rootReducer } from "./reducers";
import { actions as sqlStatsActions } from "./sqlStats";

describe("rootReducer", () => {
  it("resets redux state on RESET_STATE action", () => {
    const store = createStore(rootReducer);
    const initState = store.getState();
    const error = new Error("oops!");
    store.dispatch(sqlStatsActions.failed(error));
    const changedState = store.getState();
    store.dispatch(rootActions.resetState());
    const resetState = store.getState();

    assert.deepEqual(initState, resetState);
    assert.notDeepEqual(
      resetState.statements.lastError,
      changedState.statements.lastError,
    );
  });
});
