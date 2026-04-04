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

import { connect } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { Dispatch } from "redux";

import { AppState, uiConfigActions } from "src/store";
import { actions as statementDiagnosticsActions } from "src/store/statementDiagnostics";
import {
  actions as localStorageActions,
  updateStmtsPageLimitAction,
  updateStmsPageReqSortAction,
} from "src/store/localStorage";
import { actions as sqlStatsActions } from "src/store/sqlStats";
import { actions as nodesActions } from "../store/nodes";
import {
  StatementsPage,
  StatementsPageDispatchProps,
  StatementsPageStateProps,
} from "./statementsPage";
import {
  selectApps,
  selectDatabases,
  selectLastReset,
  selectStatements,
  selectStatementsDataValid,
  selectStatementsDataInFlight,
  selectStatementsLastError,
  selectTotalFingerprints,
  selectColumns,
  selectTimeScale,
  selectSortSetting,
  selectFilters,
  selectSearch,
  selectStatementsLastUpdated,
} from "./statementsPage.selectors";
import {
  selectStmtsPageLimit,
  selectStmtsPageReqSort,
} from "../store/utils/selectors";
import {
  selectIsTenant,
  selectHasViewActivityRedactedRole,
  selectHasAdminRole,
} from "../store/uiConfig";
import { nodeRegionsByIDSelector } from "../store/nodes";
import { StatementsRequest } from "src/api/statementsApi";
import { TimeScale } from "../timeScaleDropdown";
import { cockroach, google } from "@cockroachlabs/crdb-protobuf-client";
import { SqlStatsSortType } from "../api";

type IStatementDiagnosticsReport = cockroach.server.serverpb.IStatementDiagnosticsReport;
type IDuration = google.protobuf.IDuration;

const CreateStatementDiagnosticsReportRequest =
  cockroach.server.serverpb.CreateStatementDiagnosticsReportRequest;

const CancelStatementDiagnosticsReportRequest =
  cockroach.server.serverpb.CancelStatementDiagnosticsReportRequest;

export const ConnectedStatementsPage = withRouter(
  connect<
    StatementsPageStateProps,
    StatementsPageDispatchProps,
    RouteComponentProps
  >(
    (state: AppState, props): StatementsPageStateProps => ({
      apps: selectApps(state),
      columns: selectColumns(state),
      databases: selectDatabases(state),
      timeScale: selectTimeScale(state),
      filters: selectFilters(state),
      isTenant: selectIsTenant(state),
      hasViewActivityRedactedRole: selectHasViewActivityRedactedRole(state),
      hasAdminRole: selectHasAdminRole(state),
      lastReset: selectLastReset(state),
      nodeRegions: nodeRegionsByIDSelector(state),
      search: selectSearch(state),
      sortSetting: selectSortSetting(state),
      statements: selectStatements(state, props),
      isDataValid: selectStatementsDataValid(state),
      isReqInFlight: selectStatementsDataInFlight(state),
      lastUpdated: selectStatementsLastUpdated(state),
      statementsError: selectStatementsLastError(state),
      totalFingerprints: selectTotalFingerprints(state),
      limit: selectStmtsPageLimit(state),
      reqSortSetting: selectStmtsPageReqSort(state),
      stmtsTotalRuntimeSecs:
        state.adminUI?.statements?.data?.stmts_total_runtime_secs ?? 0,
    }),
    (dispatch: Dispatch) => ({
      refreshStatements: (req: StatementsRequest) =>
        dispatch(sqlStatsActions.refresh(req)),
      onTimeScaleChange: (ts: TimeScale) => {
        dispatch(
          sqlStatsActions.updateTimeScale({
            ts: ts,
          }),
        );
      },
      refreshStatementDiagnosticsRequests: () =>
        dispatch(statementDiagnosticsActions.refresh()),
      refreshNodes: () => dispatch(nodesActions.refresh()),
      refreshUserSQLRoles: () =>
        dispatch(uiConfigActions.refreshUserSQLRoles()),
      resetSQLStats: () => dispatch(sqlStatsActions.reset()),
      dismissAlertMessage: () =>
        dispatch(
          localStorageActions.update({
            key: "adminUi/showDiagnosticsModal",
            value: false,
          }),
        ),
      onActivateStatementDiagnostics: (
        statementFingerprint: string,
        minExecLatency: IDuration,
        expiresAfter: IDuration,
      ) => {
        dispatch(
          statementDiagnosticsActions.createReport(
            new CreateStatementDiagnosticsReportRequest({
              statement_fingerprint: statementFingerprint,
              min_execution_latency: minExecLatency,
              expires_after: expiresAfter,
            }),
          ),
        );
      },
      onSelectDiagnosticsReportDropdownOption: (
        report: IStatementDiagnosticsReport,
      ) => {
        if (!report.completed) {
          dispatch(
            statementDiagnosticsActions.cancelReport(
              new CancelStatementDiagnosticsReportRequest({
                request_id: report.id,
              }),
            ),
          );
        }
      },
      onSearchComplete: (query: string) => {
        dispatch(
          localStorageActions.update({
            key: "search/StatementsPage",
            value: query,
          }),
        );
      },
      onFilterChange: value => {
        dispatch(
          localStorageActions.update({
            key: "filters/StatementsPage",
            value: value,
          }),
        );
      },
      onSortingChange: (
        tableName: string,
        columnName: string,
        ascending: boolean,
      ) => {
        dispatch(
          localStorageActions.update({
            key: "sortSetting/StatementsPage",
            value: { columnTitle: columnName, ascending: ascending },
          }),
        );
      },
      // We use `null` when the value was never set and it will show all columns.
      // If the user modifies the selection and no columns are selected,
      // the function will save the value as a blank space, otherwise
      // it gets saved as `null`.
      onColumnsChange: (selectedColumns: string[]) =>
        dispatch(
          localStorageActions.update({
            key: "showColumns/StatementsPage",
            value:
              selectedColumns.length === 0 ? " " : selectedColumns.join(","),
          }),
        ),
      onChangeLimit: (limit: number) =>
        dispatch(updateStmtsPageLimitAction(limit)),
      onChangeReqSort: (sort: SqlStatsSortType) =>
        dispatch(updateStmsPageReqSortAction(sort)),
    }),
  )(StatementsPage),
);
