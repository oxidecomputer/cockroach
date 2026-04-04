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
import { createSelector } from "reselect";
import { RouteComponentProps, withRouter } from "react-router-dom";
import {
  refreshNodes,
  refreshTxns,
  refreshUserSQLRoles,
} from "src/redux/apiReducers";
import { AdminUIState } from "src/redux/state";
import { txnFingerprintIdAttr } from "src/util/constants";
import { getMatchParamByName } from "src/util/query";
import { nodeRegionsByIDSelector } from "src/redux/nodes";
import {
  reqSortSetting,
  selectData,
  selectLastError,
  limitSetting,
} from "src/views/transactions/transactionsPage";
import {
  TransactionDetailsStateProps,
  TransactionDetailsDispatchProps,
  TransactionDetailsProps,
  TransactionDetails,
} from "@cockroachlabs/cluster-ui";
import { setGlobalTimeScaleAction } from "src/redux/statements";
import { selectTimeScale } from "src/redux/timeScale";

export const selectTransaction = createSelector(
  (state: AdminUIState) => state.cachedData.transactions,
  (_state: AdminUIState, props: RouteComponentProps) => props,
  (transactionState, props) => {
    const transactions = transactionState.data?.transactions;
    if (!transactions) {
      return {
        isLoading: transactionState.inFlight,
        transaction: null,
        lastUpdated: null,
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
      lastUpdated: transactionState?.setAt?.utc(),
      isValid: transactionState.valid,
    };
  },
);

export default withRouter(
  connect<TransactionDetailsStateProps, TransactionDetailsDispatchProps>(
    (
      state: AdminUIState,
      props: TransactionDetailsProps,
    ): TransactionDetailsStateProps => {
      const {
        isLoading,
        transaction,
        lastUpdated,
        isValid,
      } = selectTransaction(state, props);
      return {
        timeScale: selectTimeScale(state),
        error: selectLastError(state),
        isTenant: false,
        nodeRegions: nodeRegionsByIDSelector(state),
        statements: selectData(state)?.statements,
        transaction: transaction,
        transactionFingerprintId: getMatchParamByName(
          props.match,
          txnFingerprintIdAttr,
        ),
        isLoading: isLoading,
        lastUpdated: lastUpdated,
        isDataValid: isValid,
        limit: limitSetting.selector(state),
        reqSortSetting: reqSortSetting.selector(state),
      };
    },
    {
      refreshData: refreshTxns,
      refreshNodes,
      refreshUserSQLRoles,
      onTimeScaleChange: setGlobalTimeScaleAction,
    },
  )(TransactionDetails),
);
