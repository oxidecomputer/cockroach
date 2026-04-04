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

import { createSelector } from "@reduxjs/toolkit";
import { connect } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { Dispatch } from "redux";

import { AppState, uiConfigActions } from "src/store";
import { actions as nodesActions } from "../store/nodes";
import { actions as sqlStatsActions } from "src/store/sqlStats";
import { actions as txnStatsActions } from "src/store/transactionStats";
import {
  TransactionDetails,
  TransactionDetailsDispatchProps,
  TransactionDetailsProps,
  TransactionDetailsStateProps,
} from "./transactionDetails";
import {
  selectTransactionsData,
  selectTransactionsLastError,
} from "../transactionsPage/transactionsPage.selectors";
import {
  selectIsTenant,
  selectHasViewActivityRedactedRole,
} from "../store/uiConfig";
import { nodeRegionsByIDSelector } from "../store/nodes";
import {
  selectTimeScale,
  selectTxnsPageLimit,
  selectTxnsPageReqSort,
} from "../store/utils/selectors";
import { StatementsRequest } from "src/api/statementsApi";
import { txnFingerprintIdAttr, getMatchParamByName } from "../util";
import { TimeScale } from "../timeScaleDropdown";

export const selectTransaction = createSelector(
  (state: AppState) => state.adminUI?.transactions,
  (_state: AppState, props: RouteComponentProps) => props,
  (transactionState, props) => {
    const transactions = transactionState.data?.transactions;
    if (!transactions) {
      return {
        isLoading: transactionState.inFlight,
        transaction: null,
        isValid: transactionState.valid,
      };
    }
    const txnFingerprintId = getMatchParamByName(
      props.match,
      txnFingerprintIdAttr,
    );

    const transaction = transactions.find(
      txn =>
        txn.stats_data.transaction_fingerprint_id.toString() ==
        txnFingerprintId,
    );

    return {
      isLoading: transactionState.inFlight,
      transaction: transaction,
      lastUpdated: transactionState.lastUpdated,
      isValid: transactionState.valid,
    };
  },
);

const mapStateToProps = (
  state: AppState,
  props: TransactionDetailsProps,
): TransactionDetailsStateProps => {
  const { isLoading, transaction, lastUpdated, isValid } = selectTransaction(
    state,
    props,
  );

  return {
    timeScale: selectTimeScale(state),
    error: selectTransactionsLastError(state),
    isTenant: selectIsTenant(state),
    nodeRegions: nodeRegionsByIDSelector(state),
    statements: selectTransactionsData(state)?.statements,
    transaction,
    transactionFingerprintId: getMatchParamByName(
      props.match,
      txnFingerprintIdAttr,
    ),
    isLoading: isLoading,
    hasViewActivityRedactedRole: selectHasViewActivityRedactedRole(state),
    isDataValid: isValid,
    lastUpdated,
    limit: selectTxnsPageLimit(state),
    reqSortSetting: selectTxnsPageReqSort(state),
  };
};

const mapDispatchToProps = (
  dispatch: Dispatch,
): TransactionDetailsDispatchProps => ({
  refreshData: (req?: StatementsRequest) =>
    dispatch(txnStatsActions.refresh(req)),
  refreshNodes: () => dispatch(nodesActions.refresh()),
  refreshUserSQLRoles: () => dispatch(uiConfigActions.refreshUserSQLRoles()),
  onTimeScaleChange: (ts: TimeScale) => {
    dispatch(
      sqlStatsActions.updateTimeScale({
        ts: ts,
      }),
    );
  },
});

export const TransactionDetailsPageConnected = withRouter<any, any>(
  connect(mapStateToProps, mapDispatchToProps)(TransactionDetails),
);
