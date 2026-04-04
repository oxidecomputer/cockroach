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
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { DOMAIN_NAME } from "../utils";
import { StatementsRequest } from "src/api/statementsApi";
import { TimeScale } from "../../timeScaleDropdown";
import moment from "moment";

export type StatementsResponse = cockroach.server.serverpb.StatementsResponse;

export type SQLStatsState = {
  data: StatementsResponse;
  lastError: Error;
  valid: boolean;
  lastUpdated: moment.Moment | null;
  inFlight: boolean;
};

const initialState: SQLStatsState = {
  data: null,
  lastError: null,
  valid: true,
  lastUpdated: null,
  inFlight: false,
};

export type UpdateTimeScalePayload = {
  ts: TimeScale;
};

// This is actually statements only, despite the SQLStatsState name.
// We can rename this in the future. Leaving it now to reduce backport surface area.
const sqlStatsSlice = createSlice({
  name: `${DOMAIN_NAME}/sqlstats`,
  initialState,
  reducers: {
    received: (state, action: PayloadAction<StatementsResponse>) => {
      state.inFlight = false;
      state.data = action.payload;
      state.valid = true;
      state.lastError = null;
      state.lastUpdated = moment.utc();
      state.inFlight = false;
    },
    failed: (state, action: PayloadAction<Error>) => {
      state.valid = false;
      state.lastError = action.payload;
      state.lastUpdated = moment.utc();
      state.inFlight = false;
    },
    invalidated: state => {
      state.inFlight = false;
      state.valid = false;
    },
    refresh: (state, _action: PayloadAction<StatementsRequest>) => {
      state.inFlight = true;
    },
    request: (state, _action: PayloadAction<StatementsRequest>) => {
      state.inFlight = true;
    },
    updateTimeScale: (_, _action: PayloadAction<UpdateTimeScalePayload>) => {},
    reset: () => {},
  },
});

export const { reducer, actions } = sqlStatsSlice;
