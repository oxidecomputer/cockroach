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
import { actions as nodesActions } from "src/store/nodes";
import { actions as sqlStatsActions } from "src/store/sqlStats";
import { TransactionsPage } from "./transactionsPage";
import { actions as txnStatsActions } from "src/store/transactionStats";
import {
  TransactionsPageStateProps,
  TransactionsPageDispatchProps,
} from "./transactionsPage";
import {
  selectTransactionsData,
  selectTransactionsLastError,
  selectTxnColumns,
  selectSortSetting,
  selectFilters,
  selectSearch,
  selectTransactionsDataValid,
  selectTransactionsLastUpdated,
  selectTransactionsDataInFlight,
} from "./transactionsPage.selectors";
import { selectHasAdminRole, selectIsTenant } from "../store/uiConfig";
import { nodeRegionsByIDSelector } from "../store/nodes";
import {
  selectTxnsPageLimit,
  selectTxnsPageReqSort,
  selectTimeScale,
} from "../store/utils/selectors";
import { SqlStatsSortType, StatementsRequest } from "src/api/statementsApi";
import {
  actions as localStorageActions,
  updateTxnsPageLimitAction,
  updateTxnsPageReqSortAction,
} from "../store/localStorage";
import { Filters } from "../queryFilter";
import { TimeScale } from "../timeScaleDropdown";

export const TransactionsPageConnected = withRouter(
  connect<
    TransactionsPageStateProps,
    TransactionsPageDispatchProps,
    RouteComponentProps
  >(
    (state: AppState) => ({
      columns: selectTxnColumns(state),
      data: selectTransactionsData(state),
      isDataValid: selectTransactionsDataValid(state),
      isReqInFlight: selectTransactionsDataInFlight(state),
      lastUpdated: selectTransactionsLastUpdated(state),
      timeScale: selectTimeScale(state),
      error: selectTransactionsLastError(state),
      filters: selectFilters(state),
      isTenant: selectIsTenant(state),
      nodeRegions: nodeRegionsByIDSelector(state),
      search: selectSearch(state),
      sortSetting: selectSortSetting(state),
      hasAdminRole: selectHasAdminRole(state),
      limit: selectTxnsPageLimit(state),
      reqSortSetting: selectTxnsPageReqSort(state),
    }),
    (dispatch: Dispatch) => ({
      refreshData: (req: StatementsRequest) =>
        dispatch(txnStatsActions.refresh(req)),
      refreshNodes: () => dispatch(nodesActions.refresh()),
      refreshUserSQLRoles: () =>
        dispatch(uiConfigActions.refreshUserSQLRoles()),
      resetSQLStats: () => dispatch(sqlStatsActions.reset()),
      onTimeScaleChange: (ts: TimeScale) => {
        dispatch(
          sqlStatsActions.updateTimeScale({
            ts: ts,
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
            key: "showColumns/TransactionPage",
            value:
              selectedColumns.length === 0 ? " " : selectedColumns.join(","),
          }),
        ),
      onSortingChange: (
        tableName: string,
        columnName: string,
        ascending: boolean,
      ) => {
        dispatch(
          localStorageActions.update({
            key: "sortSetting/TransactionsPage",
            value: { columnTitle: columnName, ascending: ascending },
          }),
        );
      },
      onFilterChange: (value: Filters) => {
        dispatch(
          localStorageActions.update({
            key: "filters/TransactionsPage",
            value: value,
          }),
        );
      },
      onSearchComplete: (query: string) => {
        dispatch(
          localStorageActions.update({
            key: "search/TransactionsPage",
            value: query,
          }),
        );
      },
      onChangeLimit: (limit: number) =>
        dispatch(updateTxnsPageLimitAction(limit)),
      onChangeReqSort: (sort: SqlStatsSortType) =>
        dispatch(updateTxnsPageReqSortAction(sort)),
    }),
  )(TransactionsPage),
);
