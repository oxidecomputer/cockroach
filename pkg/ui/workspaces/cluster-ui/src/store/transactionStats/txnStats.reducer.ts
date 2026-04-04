// Copyright 2023 The Cockroach Authors.
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
import { StatementsRequest } from "src/api/statementsApi";
import { TimeScale } from "../../timeScaleDropdown";
import moment from "moment";
import { StatementsResponse } from "../sqlStats";

export type TxnStatsState = {
  // Note that we request transactions from the
  // statements api, hence the StatementsResponse type here.
  data: StatementsResponse;
  inFlight: boolean;
  lastError: Error;
  valid: boolean;
  lastUpdated: moment.Moment | null;
};

const initialState: TxnStatsState = {
  data: null,
  inFlight: false,
  lastError: null,
  valid: false,
  lastUpdated: null,
};

const txnStatsSlice = createSlice({
  name: `${DOMAIN_NAME}/txnStats`,
  initialState,
  reducers: {
    received: (state, action: PayloadAction<StatementsResponse>) => {
      state.inFlight = false;
      state.data = action.payload;
      state.valid = true;
      state.lastError = null;
      state.lastUpdated = moment.utc();
    },
    failed: (state, action: PayloadAction<Error>) => {
      state.inFlight = false;
      state.valid = false;
      state.lastError = action.payload;
      state.lastUpdated = moment.utc();
    },
    invalidated: state => {
      state.inFlight = false;
      state.valid = false;
    },
    refresh: (state, _: PayloadAction<StatementsRequest>) => {
      state.inFlight = true;
    },
    request: (state, _: PayloadAction<StatementsRequest>) => {
      state.inFlight = true;
    },
  },
});

export const { reducer, actions } = txnStatsSlice;
