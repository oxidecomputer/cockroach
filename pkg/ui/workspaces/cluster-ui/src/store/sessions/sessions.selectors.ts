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

import { AppState } from "src/store";
import { createSelector } from "reselect";
import { RouteComponentProps } from "react-router-dom";
import { SessionsState } from "src/store/sessions";
import { sessionAttr } from "src/util/constants";
import { getMatchParamByName } from "src/util/query";
import { byteArrayToUuid } from "src/sessions/sessionsTable";

export const selectSession = createSelector(
  (state: AppState) => state.adminUI.sessions,
  (_state: AppState, props: RouteComponentProps) => props,
  (state: SessionsState, props: RouteComponentProps<any>) => {
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

export const selectSessionDetailsUiConfig = createSelector(
  (state: AppState) => state.adminUI.uiConfig.pages.sessionDetails,
  statementDetailsUiConfig => statementDetailsUiConfig,
);
