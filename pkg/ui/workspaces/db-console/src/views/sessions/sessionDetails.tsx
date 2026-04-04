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

import { getMatchParamByName } from "src/util/query";
import { sessionAttr } from "src/util/constants";
import _ from "lodash";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { createSelector } from "reselect";
import { Pick } from "src/util/pick";
import { AdminUIState } from "src/redux/state";
import { SessionsResponseMessage } from "src/util/api";
import { connect } from "react-redux";
import {
  CachedDataReducerState,
  refreshLiveness,
  refreshNodes,
  refreshSessions,
} from "src/redux/apiReducers";
import { nodeDisplayNameByIDSelector } from "src/redux/nodes";
import { SessionDetails, byteArrayToUuid } from "@cockroachlabs/cluster-ui";
import {
  terminateQueryAction,
  terminateSessionAction,
} from "src/redux/sessions/sessionsSagas";

type SessionsState = Pick<AdminUIState, "cachedData", "sessions">;

export const selectSession = createSelector(
  (state: SessionsState) => state.cachedData.sessions,
  (_state: SessionsState, props: RouteComponentProps) => props,
  (
    state: CachedDataReducerState<SessionsResponseMessage>,
    props: RouteComponentProps<any>,
  ) => {
    if (!state.data) {
      return null;
    }
    const sessionID = getMatchParamByName(props.match, sessionAttr);
    return {
      session: state.data.sessions.find(
        session => byteArrayToUuid(session.id) === sessionID,
      ),
    };
  },
);
const SessionDetailsPageConnected = withRouter(
  connect(
    (state: AdminUIState, props: RouteComponentProps) => ({
      nodeNames: nodeDisplayNameByIDSelector(state),
      session: selectSession(state, props),
      sessionError: state.cachedData.sessions.lastError,
    }),
    {
      refreshSessions,
      cancelSession: terminateSessionAction,
      cancelQuery: terminateQueryAction,
      refreshNodes: refreshNodes,
      refreshNodesLiveness: refreshLiveness,
    },
  )(SessionDetails),
);

export default SessionDetailsPageConnected;
