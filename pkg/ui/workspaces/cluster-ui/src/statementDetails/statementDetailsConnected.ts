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

import { withRouter } from "react-router-dom";
import { connect } from "react-redux";
import { Dispatch } from "redux";
import { RouteComponentProps } from "react-router-dom";
import {
  StatementDetails,
  StatementDetailsDispatchProps,
} from "./statementDetails";
import { AppState, uiConfigActions } from "../store";
import {
  selectStatementDetails,
  selectStatementDetailsUiConfig,
} from "./statementDetails.selectors";
import {
  selectIsTenant,
  selectHasViewActivityRedactedRole,
} from "../store/uiConfig";
import {
  nodeDisplayNameByIDSelector,
  nodeRegionsByIDSelector,
} from "../store/nodes";
import { actions as sqlDetailsStatsActions } from "src/store/statementDetails";
import { actions as sqlStatsActions } from "src/store/sqlStats";
import {
  actions as statementDiagnosticsActions,
  selectDiagnosticsReportsByStatementFingerprint,
} from "src/store/statementDiagnostics";
import { actions as localStorageActions } from "src/store/localStorage";
import { actions as nodesActions } from "../store/nodes";
import { actions as nodeLivenessActions } from "../store/liveness";
import { selectTimeScale } from "../statementsPage/statementsPage.selectors";
import { cockroach, google } from "@cockroachlabs/crdb-protobuf-client";
import { StatementDetailsRequest } from "../api";
import { TimeScale } from "../timeScaleDropdown";
import { getMatchParamByName, statementAttr } from "../util";
type IDuration = google.protobuf.IDuration;
type IStatementDiagnosticsReport = cockroach.server.serverpb.IStatementDiagnosticsReport;

const CreateStatementDiagnosticsReportRequest =
  cockroach.server.serverpb.CreateStatementDiagnosticsReportRequest;

const CancelStatementDiagnosticsReportRequest =
  cockroach.server.serverpb.CancelStatementDiagnosticsReportRequest;

// For tenant cases, we don't show information about node, regions and
// diagnostics.
const mapStateToProps = (state: AppState, props: RouteComponentProps) => {
  const { statementDetails, isLoading, lastError } = selectStatementDetails(
    state,
    props,
  );
  return {
    statementFingerprintID: getMatchParamByName(props.match, statementAttr),
    statementDetails,
    isLoading: isLoading,
    latestQuery: state.adminUI.sqlDetailsStats.latestQuery,
    latestFormattedQuery: state.adminUI.sqlDetailsStats.latestFormattedQuery,
    statementsError: lastError,
    timeScale: selectTimeScale(state),
    nodeNames: selectIsTenant(state) ? {} : nodeDisplayNameByIDSelector(state),
    nodeRegions: selectIsTenant(state) ? {} : nodeRegionsByIDSelector(state),
    diagnosticsReports:
      selectIsTenant(state) || selectHasViewActivityRedactedRole(state)
        ? []
        : selectDiagnosticsReportsByStatementFingerprint(
            state,
            state.adminUI.sqlDetailsStats.latestQuery,
          ),
    uiConfig: selectStatementDetailsUiConfig(state),
    isTenant: selectIsTenant(state),
    hasViewActivityRedactedRole: selectHasViewActivityRedactedRole(state),
  };
};

const mapDispatchToProps = (
  dispatch: Dispatch,
): StatementDetailsDispatchProps => ({
  refreshStatementDetails: (req: StatementDetailsRequest) =>
    dispatch(sqlDetailsStatsActions.refresh(req)),
  refreshStatementDiagnosticsRequests: () =>
    dispatch(statementDiagnosticsActions.refresh()),
  refreshNodes: () => dispatch(nodesActions.refresh()),
  refreshNodesLiveness: () => dispatch(nodeLivenessActions.refresh()),
  refreshUserSQLRoles: () => dispatch(uiConfigActions.refreshUserSQLRoles()),
  onTimeScaleChange: (ts: TimeScale) => {
    dispatch(
      sqlStatsActions.updateTimeScale({
        ts: ts,
      }),
    );
  },
  dismissStatementDiagnosticsAlertMessage: () =>
    dispatch(
      localStorageActions.update({
        key: "adminUi/showDiagnosticsModal",
        value: false,
      }),
    ),
  createStatementDiagnosticsReport: (
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
  onDiagnosticCancelRequest: (report: IStatementDiagnosticsReport) => {
    dispatch(
      statementDiagnosticsActions.cancelReport(
        new CancelStatementDiagnosticsReportRequest({
          request_id: report.id,
        }),
      ),
    );
  },
  onStatementDetailsQueryChange: (latestQuery: string) => {
    dispatch(sqlDetailsStatsActions.setLatestQuery(latestQuery));
  },
  onStatementDetailsFormattedQueryChange: (latestFormattedQuery: string) => {
    dispatch(
      sqlDetailsStatsActions.setLatestFormattedQuery(latestFormattedQuery),
    );
  },
});

export const ConnectedStatementDetailsPage = withRouter<any, any>(
  connect(mapStateToProps, mapDispatchToProps)(StatementDetails),
);
