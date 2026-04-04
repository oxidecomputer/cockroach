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

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { fetchData } from "src/api";

const SESSIONS_PATH = "/_status/sessions";

export type SessionsRequestMessage = cockroach.server.serverpb.ListSessionsRequest;
export type SessionsResponseMessage = cockroach.server.serverpb.ListSessionsResponse;

// getSessions gets all cluster sessions.
export const getSessions = (): Promise<SessionsResponseMessage> => {
  return fetchData(
    cockroach.server.serverpb.ListSessionsResponse,
    SESSIONS_PATH,
  );
};
