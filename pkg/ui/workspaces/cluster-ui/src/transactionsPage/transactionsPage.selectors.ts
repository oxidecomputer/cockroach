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

import { createSelector } from "reselect";

import { localStorageSelector } from "../store/utils/selectors";
import { txnStatsSelector } from "../store/transactionStats/txnStats.selector";

export const selectTransactionsData = createSelector(
  txnStatsSelector,
  transactionsState => transactionsState?.data,
);

export const selectTransactionsDataValid = createSelector(
  txnStatsSelector,
  state => state?.valid,
);

export const selectTransactionsDataInFlight = createSelector(
  txnStatsSelector,
  state => state?.inFlight,
);

export const selectTransactionsLastUpdated = createSelector(
  txnStatsSelector,
  state => state.lastUpdated,
);

export const selectTransactionsLastError = createSelector(
  txnStatsSelector,
  state => state.lastError,
);

export const selectTxnColumns = createSelector(
  localStorageSelector,
  // return array of columns if user have customized it or `null` otherwise
  localStorage =>
    localStorage["showColumns/TransactionPage"]
      ? localStorage["showColumns/TransactionPage"].split(",")
      : null,
);

export const selectSortSetting = createSelector(
  localStorageSelector,
  localStorage => localStorage["sortSetting/TransactionsPage"],
);

export const selectFilters = createSelector(
  localStorageSelector,
  localStorage => localStorage["filters/TransactionsPage"],
);

export const selectSearch = createSelector(
  localStorageSelector,
  localStorage => localStorage["search/TransactionsPage"],
);
