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
import { DOMAIN_NAME, noopReducer } from "../utils";
import { ICancelQueryRequest, ICancelSessionRequest } from ".";

type CancelQueryResponse = cockroach.server.serverpb.CancelQueryResponse;

export type TerminateQueryState = {
  data: CancelQueryResponse;
  lastError: Error;
  valid: boolean;
};

const initialState: TerminateQueryState = {
  data: null,
  lastError: null,
  valid: true,
};

const terminateQuery = createSlice({
  name: `${DOMAIN_NAME}/terminateQuery`,
  initialState,
  reducers: {
    terminateSession: (
      _state,
      _action: PayloadAction<ICancelSessionRequest>,
    ) => {},
    terminateSessionCompleted: noopReducer,
    terminateSessionFailed: (_state, _action: PayloadAction<Error>) => {},
    terminateQuery: (_state, _action: PayloadAction<ICancelQueryRequest>) => {},
    terminateQueryCompleted: noopReducer,
    terminateQueryFailed: (_state, _action: PayloadAction<Error>) => {},
  },
});

export const { reducer, actions } = terminateQuery;
