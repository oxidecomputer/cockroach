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

import { createMemoryHistory } from "history";
import { SessionDetailsProps } from "./sessionDetails";
import {
  activeSession,
  idleSession,
  idleTransactionSession,
} from "./sessionsPage.fixture";
import { sessionAttr } from "src/util/constants";
import {
  CancelSessionRequestMessage,
  CancelQueryRequestMessage,
} from "src/api/terminateQueryApi";

const history = createMemoryHistory({ initialEntries: ["/sessions"] });

const sessionDetailsPropsBase: SessionDetailsProps = {
  id: "blah",
  nodeNames: {
    1: "localhost",
  },
  sessionError: null,
  session: null,
  history,
  location: {
    pathname: "/sessions/blah",
    search: "",
    hash: "",
    state: null,
  },
  match: {
    path: "/sessions/blah",
    url: "/sessions/blah",
    isExact: true,
    params: { [sessionAttr]: "blah" },
  },

  refreshSessions: () => {},
  cancelSession: (req: CancelSessionRequestMessage) => {},
  cancelQuery: (req: CancelQueryRequestMessage) => {},
  refreshNodes: () => {},
  refreshNodesLiveness: () => {},
  uiConfig: {
    showGatewayNodeLink: true,
  },
};

export const sessionDetailsIdlePropsFixture: SessionDetailsProps = {
  ...sessionDetailsPropsBase,
  session: idleSession,
};

export const sessionDetailsActiveTxnPropsFixture: SessionDetailsProps = {
  ...sessionDetailsPropsBase,
  session: idleTransactionSession,
};

export const sessionDetailsActiveStmtPropsFixture: SessionDetailsProps = {
  ...sessionDetailsPropsBase,
  session: activeSession,
};

export const sessionDetailsNotFound: SessionDetailsProps = {
  ...sessionDetailsPropsBase,
  session: { session: null },
};
