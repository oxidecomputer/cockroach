// Copyright 2022 The Cockroach Authors.
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
import { DOMAIN_NAME } from "../utils";
import {
  ErrorWithKey,
  StatementDetailsRequest,
  StatementDetailsResponse,
  StatementDetailsResponseWithKey,
} from "src/api/statementsApi";
import { generateStmtDetailsToID } from "../../util";

export type SQLDetailsStatsState = {
  data: StatementDetailsResponse;
  lastError: Error;
  valid: boolean;
  inFlight: boolean;
};

export type SQLDetailsStatsReducerState = {
  cachedData: {
    [id: string]: SQLDetailsStatsState;
  };
  latestQuery: string;
  latestFormattedQuery: string;
};

const initialState: SQLDetailsStatsReducerState = {
  cachedData: {},
  latestQuery: "",
  latestFormattedQuery: "",
};

const sqlDetailsStatsSlice = createSlice({
  name: `${DOMAIN_NAME}/sqlDetailsStats`,
  initialState,
  reducers: {
    received: (
      state,
      action: PayloadAction<StatementDetailsResponseWithKey>,
    ) => {
      state.cachedData[action.payload.key] = {
        data: action.payload.stmtResponse,
        valid: true,
        lastError: null,
        inFlight: false,
      };
    },
    failed: (state, action: PayloadAction<ErrorWithKey>) => {
      state.cachedData[action.payload.key] = {
        data: null,
        valid: false,
        lastError: action.payload.err,
        inFlight: false,
      };
    },
    invalidated: (state, action: PayloadAction<{ key: string }>) => {
      delete state.cachedData[action.payload.key];
    },
    invalidateAll: state => {
      const keys = Object.keys(state);
      for (const key in keys) {
        delete state.cachedData[key];
      }
    },
    refresh: (state, action: PayloadAction<StatementDetailsRequest>) => {
      const key = action?.payload
        ? generateStmtDetailsToID(
            action.payload.fingerprint_id,
            action.payload.app_names.toString(),
            action.payload.start,
            action.payload.end,
          )
        : "";
      state.cachedData[key] = {
        data: null,
        valid: false,
        lastError: null,
        inFlight: true,
      };
    },
    request: (state, action: PayloadAction<StatementDetailsRequest>) => {
      const key = action?.payload
        ? generateStmtDetailsToID(
            action.payload.fingerprint_id,
            action.payload.app_names.toString(),
            action.payload.start,
            action.payload.end,
          )
        : "";
      state.cachedData[key] = {
        data: null,
        valid: false,
        lastError: null,
        inFlight: true,
      };
    },
    setLatestQuery: (state, action: PayloadAction<string>) => {
      state.latestQuery = action.payload;
    },
    setLatestFormattedQuery: (state, action: PayloadAction<string>) => {
      state.latestFormattedQuery = action.payload;
    },
  },
});

export const { reducer, actions } = sqlDetailsStatsSlice;
