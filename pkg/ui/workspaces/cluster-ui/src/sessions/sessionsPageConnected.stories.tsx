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

import React from "react";
import { storiesOf } from "@storybook/react";
import { createMemoryHistory } from "history";
import createSagaMiddleware from "redux-saga";
import { Provider } from "react-redux";
import {
  ConnectedRouter,
  connectRouter,
  routerMiddleware,
} from "connected-react-router";
import {
  applyMiddleware,
  combineReducers,
  compose,
  createStore,
  Store,
} from "redux";
import { AppState, rootReducer, sagas } from "../store";

const history = createMemoryHistory();
const routerReducer = connectRouter(history);
const sagaMiddleware = createSagaMiddleware();

const store: Store<AppState> = createStore(
  combineReducers({
    router: routerReducer,
    adminUI: rootReducer,
  }),
  compose(
    applyMiddleware(sagaMiddleware, routerMiddleware(history)),
    (window as any).__REDUX_DEVTOOLS_EXTENSION__ &&
      (window as any).__REDUX_DEVTOOLS_EXTENSION__(),
  ),
);

sagaMiddleware.run(sagas);

storiesOf("Sessions Page Connected", module)
  .addDecorator(storyFn => (
    <Provider store={store}>
      <ConnectedRouter history={history}>{storyFn()}</ConnectedRouter>
    </Provider>
  ))
  .addDecorator(storyFn => (
    <div style={{ backgroundColor: "#F5F7FA" }}>{storyFn()}</div>
  ));
