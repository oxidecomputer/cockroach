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

import { isUndefined } from "lodash";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { DiagnosticStatuses } from "src/statementsDiagnostics";

type IStatementDiagnosticsReport = cockroach.server.serverpb.IStatementDiagnosticsReport;

export function getDiagnosticsStatus(
  diagnosticsRequest: IStatementDiagnosticsReport,
): DiagnosticStatuses {
  if (diagnosticsRequest.completed) {
    return "READY";
  }

  return "WAITING";
}

export function sortByRequestedAtField(
  a: IStatementDiagnosticsReport,
  b: IStatementDiagnosticsReport,
) {
  const activatedOnA = a.requested_at?.seconds?.toNumber();
  const activatedOnB = b.requested_at?.seconds?.toNumber();
  if (isUndefined(activatedOnA) && isUndefined(activatedOnB)) {
    return 0;
  }
  if (activatedOnA < activatedOnB) {
    return -1;
  }
  if (activatedOnA > activatedOnB) {
    return 1;
  }
  return 0;
}

export function sortByCompletedField(
  a: IStatementDiagnosticsReport,
  b: IStatementDiagnosticsReport,
) {
  const completedA = a.completed ? 1 : -1;
  const completedB = b.completed ? 1 : -1;
  if (completedA < completedB) {
    return -1;
  }
  if (completedA > completedB) {
    return 1;
  }
  return 0;
}

export function sortByStatementFingerprintField(
  a: IStatementDiagnosticsReport,
  b: IStatementDiagnosticsReport,
) {
  const statementFingerprintA = a.statement_fingerprint;
  const statementFingerprintB = b.statement_fingerprint;
  if (
    isUndefined(statementFingerprintA) &&
    isUndefined(statementFingerprintB)
  ) {
    return 0;
  }
  if (statementFingerprintA < statementFingerprintB) {
    return -1;
  }
  if (statementFingerprintA > statementFingerprintB) {
    return 1;
  }
  return 0;
}
