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

import { connect } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { AppState } from "src/store";
import { SessionDetails } from ".";
import {
  actions as sessionsActions,
  selectSession,
  selectSessionDetailsUiConfig,
} from "src/store/sessions";
import { actions as terminateQueryActions } from "src/store/terminateQuery";
import { actions as nodesActions } from "src/store/nodes";
import { actions as nodesLivenessActions } from "src/store/liveness";

import { nodeDisplayNameByIDSelector } from "src/store/nodes";
import { selectIsTenant } from "../store/uiConfig";

export const SessionDetailsPageConnected = withRouter(
  connect(
    (state: AppState, props: RouteComponentProps) => ({
      nodeNames: nodeDisplayNameByIDSelector(state),
      session: selectSession(state, props),
      sessionError: state.adminUI.sessions.lastError,
      uiConfig: selectSessionDetailsUiConfig(state),
      isTenant: selectIsTenant(state),
    }),
    {
      refreshSessions: sessionsActions.refresh,
      cancelSession: terminateQueryActions.terminateSession,
      cancelQuery: terminateQueryActions.terminateQuery,
      refreshNodes: nodesActions.refresh,
      refreshNodesLiveness: nodesLivenessActions.refresh,
    },
  )(SessionDetails),
);
