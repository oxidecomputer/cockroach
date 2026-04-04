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

import { assert } from "chai";
import { selectDiagnosticsReportsPerStatement } from "./statementDiagnostics.selectors";
import Long from "long";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import IStatementDiagnosticsReport = cockroach.server.serverpb.IStatementDiagnosticsReport;

const reports: IStatementDiagnosticsReport[] = [
  {
    id: Long.fromNumber(1),
    completed: false,
    statement_fingerprint: "SHOW database",
    statement_diagnostics_id: Long.fromNumber(594413981435920385),
    requested_at: { seconds: Long.fromNumber(100), nanos: 737251000 },
  },
  {
    id: Long.fromNumber(2),
    completed: true,
    statement_fingerprint: "SHOW database",
    statement_diagnostics_id: Long.fromNumber(594413281435920385),
    requested_at: { seconds: Long.fromNumber(200), nanos: 737251000 },
  },
  {
    id: Long.fromNumber(3),
    completed: true,
    statement_fingerprint: "SHOW database",
    statement_diagnostics_id: Long.fromNumber(594413281435920385),
    requested_at: { seconds: Long.fromNumber(300), nanos: 737251000 },
  },
];

describe("statementDiagnostics selectors", () => {
  describe("selectDiagnosticsReportsPerStatement", () => {
    it("returns diagnostics reports sorted in descending order", () => {
      const diagnosticsPerStatement = selectDiagnosticsReportsPerStatement.resultFunc(
        reports,
      );
      assert.deepEqual(
        diagnosticsPerStatement["SHOW database"][0].id,
        Long.fromNumber(3),
      );
    });
  });
});
