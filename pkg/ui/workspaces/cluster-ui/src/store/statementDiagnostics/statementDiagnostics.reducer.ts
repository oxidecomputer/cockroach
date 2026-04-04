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

type CreateStatementDiagnosticsReportRequest = cockroach.server.serverpb.CreateStatementDiagnosticsReportRequest;
type CancelStatementDiagnosticsReportRequest = cockroach.server.serverpb.CancelStatementDiagnosticsReportRequest;
type StatementDiagnosticsReportsResponse = cockroach.server.serverpb.StatementDiagnosticsReportsResponse;

export type StatementDiagnosticsState = {
  data: StatementDiagnosticsReportsResponse;
  lastError: Error;
  valid: boolean;
};

const initialState: StatementDiagnosticsState = {
  data: null,
  valid: true,
  lastError: null,
};

const statementDiagnosticsSlice = createSlice({
  name: `${DOMAIN_NAME}/statementDiagnostics`,
  initialState,
  reducers: {
    received: (
      state: StatementDiagnosticsState,
      action: PayloadAction<StatementDiagnosticsReportsResponse>,
    ) => {
      state.data = action.payload;
      state.lastError = null;
      state.valid = true;
    },
    failed: (
      state: StatementDiagnosticsState,
      action: PayloadAction<Error>,
    ) => {
      state.lastError = action.payload;
      state.valid = false;
    },
    refresh: noopReducer,
    request: noopReducer,
    invalidated: noopReducer,
    createReport: (
      _state,
      _action: PayloadAction<CreateStatementDiagnosticsReportRequest>,
    ) => {},
    createReportCompleted: noopReducer,
    createReportFailed: (_state, _action: PayloadAction<Error>) => {},
    cancelReport: (
      _state,
      _action: PayloadAction<CancelStatementDiagnosticsReportRequest>,
    ) => {},
    cancelReportCompleted: noopReducer,
    cancelReportFailed: (_state, _action: PayloadAction<Error>) => {},
    openNewDiagnosticsModal: (_state, _action: PayloadAction<string>) => {},
  },
});

export const { actions, reducer } = statementDiagnosticsSlice;
