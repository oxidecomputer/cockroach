// Copyright 2018 The Cockroach Authors.
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
import { Route, Redirect } from "react-router-dom";
import { Store } from "redux";

import { doLogout, selectLoginState } from "src/redux/login";
import { AdminUIState, AppDispatch } from "src/redux/state";
import LoginPage from "src/views/login/loginPage";

export const LOGIN_PAGE = "/login";
export const LOGOUT_PAGE = "/logout";

export function createLoginRoute() {
  return <Route exact path={LOGIN_PAGE} component={LoginPage} />;
}

export function createLogoutRoute(store: Store<AdminUIState>): JSX.Element {
  return (
    <Route
      exact
      path={LOGOUT_PAGE}
      render={() => {
        const loginState = selectLoginState(store.getState());

        if (!loginState.loggedInUser()) {
          return <Redirect to={LOGIN_PAGE} />;
        }

        (store.dispatch as AppDispatch)(doLogout());
      }}
    />
  );
}
