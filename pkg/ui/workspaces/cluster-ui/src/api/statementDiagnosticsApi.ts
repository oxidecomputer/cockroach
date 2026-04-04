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

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { fetchData } from "src/api";

const STATEMENT_DIAGNOSTICS_PATH = "/_status/stmtdiagreports";
const CREATE_STATEMENT_DIAGNOSTICS_REPORT_PATH = "/_status/stmtdiagreports";
const CANCEL_STATEMENT_DIAGNOSTICS_REPORT_PATH =
  "/_status/stmtdiagreports/cancel";

type CreateStatementDiagnosticsReportRequestMessage = cockroach.server.serverpb.CreateStatementDiagnosticsReportRequest;
type CreateStatementDiagnosticsReportResponseMessage = cockroach.server.serverpb.CreateStatementDiagnosticsReportResponse;
type CancelStatementDiagnosticsReportRequestMessage = cockroach.server.serverpb.CancelStatementDiagnosticsReportRequest;
type CancelStatementDiagnosticsReportResponseMessage = cockroach.server.serverpb.CancelStatementDiagnosticsReportResponse;

export function getStatementDiagnosticsReports(): Promise<
  cockroach.server.serverpb.StatementDiagnosticsReportsResponse
> {
  return fetchData(
    cockroach.server.serverpb.StatementDiagnosticsReportsResponse,
    STATEMENT_DIAGNOSTICS_PATH,
  );
}

export function createStatementDiagnosticsReport(
  req: CreateStatementDiagnosticsReportRequestMessage,
): Promise<CreateStatementDiagnosticsReportResponseMessage> {
  return fetchData(
    cockroach.server.serverpb.CreateStatementDiagnosticsReportResponse,
    CREATE_STATEMENT_DIAGNOSTICS_REPORT_PATH,
    cockroach.server.serverpb.CreateStatementDiagnosticsReportRequest,
    req,
  );
}

export function cancelStatementDiagnosticsReport(
  req: CancelStatementDiagnosticsReportRequestMessage,
): Promise<CancelStatementDiagnosticsReportResponseMessage> {
  return fetchData(
    cockroach.server.serverpb.CancelStatementDiagnosticsReportResponse,
    CANCEL_STATEMENT_DIAGNOSTICS_REPORT_PATH,
    cockroach.server.serverpb.CancelStatementDiagnosticsReportRequest,
    req,
  );
}
