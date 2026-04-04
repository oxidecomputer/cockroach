// Copyright 2018 The Cockroach Authors.
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
import { withRouter } from "react-router-dom";
import { createSelector } from "reselect";
import Long from "long";

import {
  refreshLiveness,
  refreshNodes,
  refreshStatementDiagnosticsRequests,
  refreshStatementDetails,
  refreshUserSQLRoles,
} from "src/redux/apiReducers";
import { RouteComponentProps } from "react-router";
import {
  nodeDisplayNameByIDSelector,
  nodeRegionsByIDSelector,
} from "src/redux/nodes";
import { AdminUIState, AppDispatch } from "src/redux/state";
import { selectDiagnosticsReportsByStatementFingerprint } from "src/redux/statements/statementsSelectors";
import {
  StatementDetails,
  StatementDetailsDispatchProps,
  StatementDetailsStateProps,
  toRoundedDateRange,
  util,
} from "@cockroachlabs/cluster-ui";
import {
  cancelStatementDiagnosticsReportAction,
  createStatementDiagnosticsReportAction,
  setGlobalTimeScaleAction,
} from "src/redux/statements";
import { createStatementDiagnosticsAlertLocalSetting } from "src/redux/alerts";
import { selectHasViewActivityRedactedRole } from "src/redux/user";
import * as protos from "src/js/protos";
import { StatementDetailsResponseMessage } from "src/util/api";
import { getMatchParamByName, queryByName } from "src/util/query";

import { appNamesAttr, statementAttr } from "src/util/constants";
import {
  statementDetailsLatestQueryAction,
  statementDetailsLatestFormattedQueryAction,
} from "src/redux/sqlActivity";
import { selectTimeScale } from "src/redux/timeScale";

type IStatementDiagnosticsReport = protos.cockroach.server.serverpb.IStatementDiagnosticsReport;

const { generateStmtDetailsToID } = util;

export const selectStatementDetails = createSelector(
  (_state: AdminUIState, props: RouteComponentProps): string =>
    getMatchParamByName(props.match, statementAttr),
  (_state: AdminUIState, props: RouteComponentProps): string =>
    queryByName(props.location, appNamesAttr),
  selectTimeScale,
  (state: AdminUIState) => state.cachedData.statementDetails,
  (
    fingerprintID,
    appNames,
    timeScale,
    statementDetailsStats,
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
    if (Object.keys(statementDetailsStats).includes(key)) {
      return {
        statementDetails: statementDetailsStats[key].data,
        isLoading: statementDetailsStats[key].inFlight,
        lastError: statementDetailsStats[key].lastError,
      };
    }
    return { statementDetails: null, isLoading: true, lastError: null };
  },
);

const mapStateToProps = (
  state: AdminUIState,
  props: RouteComponentProps,
): StatementDetailsStateProps => {
  const { statementDetails, isLoading, lastError } = selectStatementDetails(
    state,
    props,
  );
  return {
    statementFingerprintID: getMatchParamByName(props.match, statementAttr),
    statementDetails,
    isLoading: isLoading,
    latestQuery: state.sqlActivity.statementDetailsLatestQuery,
    latestFormattedQuery:
      state.sqlActivity.statementDetailsLatestFormattedQuery,
    statementsError: lastError,
    timeScale: selectTimeScale(state),
    nodeNames: nodeDisplayNameByIDSelector(state),
    nodeRegions: nodeRegionsByIDSelector(state),
    diagnosticsReports: selectDiagnosticsReportsByStatementFingerprint(
      state,
      state.sqlActivity.statementDetailsLatestQuery,
    ),
    hasViewActivityRedactedRole: selectHasViewActivityRedactedRole(state),
  };
};

const mapDispatchToProps: StatementDetailsDispatchProps = {
  refreshStatementDetails,
  refreshStatementDiagnosticsRequests,
  dismissStatementDiagnosticsAlertMessage: () =>
    createStatementDiagnosticsAlertLocalSetting.set({ show: false }),
  createStatementDiagnosticsReport: createStatementDiagnosticsReportAction,
  onTimeScaleChange: setGlobalTimeScaleAction,
  onDiagnosticCancelRequest: (report: IStatementDiagnosticsReport) => {
    return (dispatch: AppDispatch) => {
      dispatch(cancelStatementDiagnosticsReportAction(report.id));
    };
  },
  onStatementDetailsQueryChange: statementDetailsLatestQueryAction,
  onStatementDetailsFormattedQueryChange: statementDetailsLatestFormattedQueryAction,
  refreshNodes: refreshNodes,
  refreshNodesLiveness: refreshLiveness,
  refreshUserSQLRoles: refreshUserSQLRoles,
};

export default withRouter(
  connect<StatementDetailsStateProps, StatementDetailsDispatchProps>(
    mapStateToProps,
    mapDispatchToProps,
  )(StatementDetails),
);
