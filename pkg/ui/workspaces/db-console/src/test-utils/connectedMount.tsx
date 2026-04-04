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

import React from "react";
import { mount, ReactWrapper } from "enzyme";
import { Action, Store } from "redux";
import { Provider } from "react-redux";
import { ConnectedRouter } from "connected-react-router";
import { createMemoryHistory } from "history";

import "src/enzymeInit";
import { AdminUIState, createAdminUIStore } from "src/redux/state";

export function connectedMount(
  nodeFactory: (store: Store<AdminUIState>) => React.ReactNode,
) {
  const history = createMemoryHistory({
    initialEntries: ["/"],
  });
  const store: Store<AdminUIState, Action> = createAdminUIStore(history);
  const wrapper: ReactWrapper = mount(
    <Provider store={store}>
      <ConnectedRouter history={history}>{nodeFactory(store)}</ConnectedRouter>
    </Provider>,
  );
  return wrapper;
}
