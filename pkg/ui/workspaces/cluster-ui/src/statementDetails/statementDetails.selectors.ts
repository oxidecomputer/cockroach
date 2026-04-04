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

import Long from "long";
import { createSelector } from "@reduxjs/toolkit";
import { RouteComponentProps } from "react-router-dom";
import { AppState } from "../store";
import {
  appNamesAttr,
  statementAttr,
  getMatchParamByName,
  queryByName,
  generateStmtDetailsToID,
} from "../util";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { TimeScale, toRoundedDateRange } from "../timeScaleDropdown";
import { selectTimeScale } from "../statementsPage/statementsPage.selectors";
type StatementDetailsResponseMessage = cockroach.server.serverpb.StatementDetailsResponse;

export const selectStatementDetails = createSelector(
  (_state: AppState, props: RouteComponentProps): string =>
    getMatchParamByName(props.match, statementAttr),
  (_state: AppState, props: RouteComponentProps): string =>
    queryByName(props.location, appNamesAttr),
  (state: AppState): TimeScale => selectTimeScale(state),
  (state: AppState) => state.adminUI.sqlDetailsStats.cachedData,
  (
    fingerprintID,
    appNames,
    timeScale,
    statementDetailsStatsData,
  ): {
    statementDetails: StatementDetailsResponseMessage;
    isLoading: boolean;
    lastError: Error;
  } => {
    // Since the aggregation interval is 1h, we want to round the selected timeScale to include
    // the full hour. If a timeScale is between 14:32 - 15:17 we want to search for values
    // between 14:00 - 16:00. We don't encourage the aggregation interval to be modified, but
    // in case that changes in the future we might consider changing this function to use the
    // cluster settings value for the rounding function.
    const [start, end] = toRoundedDateRange(timeScale);
    const key = generateStmtDetailsToID(
      fingerprintID,
      appNames,
      Long.fromNumber(start.unix()),
      Long.fromNumber(end.unix()),
    );
    if (Object.keys(statementDetailsStatsData).includes(key)) {
      return {
        statementDetails: statementDetailsStatsData[key].data,
        isLoading: statementDetailsStatsData[key].inFlight,
        lastError: statementDetailsStatsData[key].lastError,
      };
    }
    return { statementDetails: null, isLoading: true, lastError: null };
  },
);

export const selectStatementDetailsUiConfig = createSelector(
  (state: AppState) => state.adminUI.uiConfig.pages.statementDetails,
  statementDetailsUiConfig => statementDetailsUiConfig,
);
