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

import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { merge } from "lodash";
import { DOMAIN_NAME, noopReducer } from "../utils";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
export type UserSQLRolesRequest = cockroach.server.serverpb.UserSQLRolesRequest;

export type UIConfigState = {
  isTenant: boolean;
  userSQLRoles: string[];
  hasViewActivityRedactedRole: boolean;
  hasAdminRole: boolean;
  pages: {
    statementDetails: {
      showStatementDiagnosticsLink: boolean;
    };
    sessionDetails: {
      showGatewayNodeLink: boolean;
    };
  };
};

const initialState: UIConfigState = {
  isTenant: false,
  userSQLRoles: [],
  hasViewActivityRedactedRole: false,
  hasAdminRole: false,
  pages: {
    statementDetails: {
      showStatementDiagnosticsLink: true,
    },
    sessionDetails: {
      showGatewayNodeLink: false,
    },
  },
};

/**
 * `uiConfigSlice` is responsible to store configuration parameters which works as feature flags
 * and can be set dynamically by dispatching `update` action with updated configuration.
 * This might be useful in case client application that integrates some components or pages from
 * `cluster-ui` and has to exclude or add some extra logic on a page.
 **/
const uiConfigSlice = createSlice({
  name: `${DOMAIN_NAME}/uiConfig`,
  initialState,
  reducers: {
    update: (state, action: PayloadAction<Partial<UIConfigState>>) => {
      merge(state, action.payload);
    },
    receivedUserSQLRoles: (state, action: PayloadAction<string[]>) => {
      if (action?.payload) {
        state.userSQLRoles = action.payload;
      }
    },
    invalidatedUserSQLRoles: state => {
      state.userSQLRoles = [];
    },
    // Define actions that don't change state
    refreshUserSQLRoles: noopReducer,
    requestUserSQLRoles: noopReducer,
  },
});

export const { actions, reducer } = uiConfigSlice;
