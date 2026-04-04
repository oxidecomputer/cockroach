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

import { duration } from "moment";

export const aggregatedTsAttr = "aggregated_ts";
export const appAttr = "app";
export const appNamesAttr = "appNames";
export const ascendingAttr = "ascending";
export const columnTitleAttr = "columnTitle";
export const dashQueryString = "dash";
export const dashboardNameAttr = "dashboard_name";
export const databaseAttr = "database";
export const databaseNameAttr = "database_name";
export const fingerprintIDAttr = "fingerprint_id";
export const implicitTxnAttr = "implicitTxn";
export const nodeIDAttr = "node_id";
export const nodeQueryString = "node";
export const rangeIDAttr = "range_id";
export const statementAttr = "statement";
export const sessionAttr = "session";
export const tabAttr = "tab";
export const tableNameAttr = "table_name";
export const txnFingerprintIdAttr = "txn_fingerprint_id";
export const unset = "(unset)";
export const viewAttr = "view";

export const REMOTE_DEBUGGING_ERROR_TEXT =
  "This information is not available due to the current value of the 'server.remote_debugging.mode' setting.";

export const serverToClientErrorMessageMap = new Map([
  [
    "not allowed (due to the 'server.remote_debugging.mode' setting)",
    REMOTE_DEBUGGING_ERROR_TEXT,
  ],
]);
