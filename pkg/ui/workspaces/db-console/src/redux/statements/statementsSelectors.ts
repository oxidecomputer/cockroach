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

import { chain, orderBy } from "lodash";
import { createSelector } from "reselect";
import { AdminUIState } from "src/redux/state";
import { cockroach } from "src/js/protos";
import IStatementDiagnosticsReport = cockroach.server.serverpb.IStatementDiagnosticsReport;

export const selectStatementByFingerprint = createSelector(
  (state: AdminUIState) => state.cachedData.statements.data?.statements,
  (_state: AdminUIState, statementFingerprint: string) => statementFingerprint,
  (statements, statementFingerprint) =>
    (statements || []).find(
      statement => statement.key.key_data.query === statementFingerprint,
    ),
);

export const selectDiagnosticsReportsByStatementFingerprint = createSelector(
  (state: AdminUIState) =>
    state.cachedData.statementDiagnosticsReports.data?.reports || [],
  (_state: AdminUIState, statementFingerprint: string) => statementFingerprint,
  (requests, statementFingerprint) =>
    (requests || []).filter(
      request => request.statement_fingerprint === statementFingerprint,
    ),
);

export const selectDiagnosticsReportsCountByStatementFingerprint = createSelector(
  selectDiagnosticsReportsByStatementFingerprint,
  requests => requests.length,
);

export const selectStatementDiagnosticsReports = createSelector(
  (state: AdminUIState) =>
    state.cachedData.statementDiagnosticsReports.data?.reports,
  diagnosticsReports => diagnosticsReports,
);

export const statementDiagnosticsReportsInFlight = createSelector(
  (state: AdminUIState) =>
    state.cachedData.statementDiagnosticsReports.inFlight,
  inFlight => inFlight,
);

type StatementDiagnosticsDictionary = {
  [statementFingerprint: string]: IStatementDiagnosticsReport[];
};

export const selectDiagnosticsReportsPerStatement = createSelector(
  selectStatementDiagnosticsReports,
  (
    diagnosticsReports: IStatementDiagnosticsReport[],
  ): StatementDiagnosticsDictionary =>
    chain(diagnosticsReports)
      .groupBy(diagnosticsReport => diagnosticsReport.statement_fingerprint)
      .mapValues(diagnostics =>
        orderBy(diagnostics, d => d.requested_at.seconds.toNumber(), ["desc"]),
      )
      .value(),
);
