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

import { RouteComponentProps, withRouter } from "react-router-dom";
import { connect } from "react-redux";
import { AppState } from "src/store";
import { SessionsState } from "src/store/sessions";

import { createSelector } from "reselect";
import { OwnProps, SessionsPage } from "./index";

import { actions as sessionsActions } from "src/store/sessions";
import { actions as localStorageActions } from "src/store/localStorage";
import {
  actions as terminateQueryActions,
  ICancelQueryRequest,
  ICancelSessionRequest,
} from "src/store/terminateQuery";
import { Dispatch } from "redux";
import { Filters } from "../queryFilter";
import { sqlStatsSelector } from "../store/sqlStats/sqlStats.selector";

export const selectSessionsData = createSelector(
  sqlStatsSelector,
  sessionsState => (sessionsState.valid ? sessionsState.data : null),
);

export const adminUISelector = createSelector(
  (state: AppState) => state.adminUI,
  adminUiState => adminUiState,
);

export const localStorageSelector = createSelector(
  adminUISelector,
  adminUiState => adminUiState.localStorage,
);

export const selectSessions = createSelector(
  (state: AppState) => state.adminUI.sessions,
  (state: SessionsState) => {
    if (!state.data) {
      return null;
    }
    return state.data.sessions.map(session => {
      return { session };
    });
  },
);

export const selectAppName = createSelector(
  (state: AppState) => state.adminUI.sessions,
  (state: SessionsState) => {
    if (!state.data) {
      return null;
    }
    return state.data.internal_app_name_prefix;
  },
);

export const selectSortSetting = createSelector(
  (state: AppState) => state.adminUI.localStorage,
  localStorage => localStorage["sortSetting/SessionsPage"],
);

export const selectColumns = createSelector(
  localStorageSelector,
  localStorage =>
    localStorage["showColumns/SessionsPage"]
      ? localStorage["showColumns/SessionsPage"].split(",")
      : null,
);

export const selectFilters = createSelector(
  localStorageSelector,
  localStorage => localStorage["filters/SessionsPage"],
);

export const SessionsPageConnected = withRouter(
  connect(
    (state: AppState, props: RouteComponentProps) => ({
      sessions: selectSessions(state),
      internalAppNamePrefix: selectAppName(state),
      sessionsError: state.adminUI.sessions.lastError,
      sortSetting: selectSortSetting(state),
      columns: selectColumns(state),
      filters: selectFilters(state),
    }),
    (dispatch: Dispatch) => ({
      refreshSessions: () => dispatch(sessionsActions.refresh()),
      cancelSession: (payload: ICancelSessionRequest) =>
        dispatch(terminateQueryActions.terminateSession(payload)),
      cancelQuery: (payload: ICancelQueryRequest) =>
        dispatch(terminateQueryActions.terminateQuery(payload)),
      onSortingChange: (
        tableName: string,
        columnName: string,
        ascending: boolean,
      ) => {
        dispatch(
          localStorageActions.update({
            key: "sortSetting/SessionsPage",
            value: { columnTitle: columnName, ascending: ascending },
          }),
        );
      },
      onFilterChange: (value: Filters) => {
        dispatch(
          localStorageActions.update({
            key: "filters/SessionsPage",
            value: value,
          }),
        );
      },
      onColumnsChange: (selectedColumns: string[]) =>
        dispatch(
          localStorageActions.update({
            key: "showColumns/SessionsPage",
            value:
              selectedColumns.length === 0 ? " " : selectedColumns.join(","),
          }),
        ),
    }),
  )(SessionsPage),
);
