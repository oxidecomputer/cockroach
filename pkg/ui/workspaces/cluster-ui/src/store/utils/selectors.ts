// Copyright 2022 The Cockroach Authors.
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
import { LocalStorageKeys } from "../localStorage";
import { AppState } from "../reducers";

export const adminUISelector = createSelector(
  (state: AppState) => state.adminUI,
  adminUiState => adminUiState,
);

export const localStorageSelector = createSelector(
  adminUISelector,
  adminUiState => adminUiState?.localStorage,
);

export const selectTimeScale = createSelector(
  localStorageSelector,
  localStorage => localStorage[LocalStorageKeys.GLOBAL_TIME_SCALE],
);

export const selectStmtsPageLimit = createSelector(
  localStorageSelector,
  localStorage => localStorage[LocalStorageKeys.STMT_FINGERPRINTS_LIMIT],
);

export const selectStmtsPageReqSort = createSelector(
  localStorageSelector,
  localStorage => localStorage[LocalStorageKeys.STMT_FINGERPRINTS_SORT],
);

export const selectTxnsPageLimit = createSelector(
  localStorageSelector,
  localStorage => localStorage[LocalStorageKeys.TXN_FINGERPRINTS_LIMIT],
);

export const selectTxnsPageReqSort = createSelector(
  localStorageSelector,
  localStorage => localStorage[LocalStorageKeys.TXN_FINGERPRINTS_SORT],
);
