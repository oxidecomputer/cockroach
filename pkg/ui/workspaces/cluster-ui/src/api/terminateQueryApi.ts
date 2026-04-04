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

const STATUS_PREFIX = "/_status";

export type CancelSessionRequestMessage = cockroach.server.serverpb.CancelSessionRequest;
export type CancelSessionResponseMessage = cockroach.server.serverpb.CancelSessionResponse;
export type CancelQueryRequestMessage = cockroach.server.serverpb.CancelQueryRequest;
export type CancelQueryResponseMessage = cockroach.server.serverpb.CancelQueryResponse;

export const terminateSession = (
  req: CancelSessionRequestMessage,
): Promise<CancelSessionResponseMessage> => {
  return fetchData(
    cockroach.server.serverpb.CancelSessionResponse,
    `${STATUS_PREFIX}/cancel_session/${req.node_id}`,
    cockroach.server.serverpb.CancelSessionRequest,
    req,
  );
};

export const terminateQuery = (
  req: CancelQueryRequestMessage,
): Promise<CancelQueryResponseMessage> => {
  return fetchData(
    cockroach.server.serverpb.CancelQueryResponse,
    `${STATUS_PREFIX}/cancel_query/${req.node_id}`,
    cockroach.server.serverpb.CancelQueryRequest,
    req,
  );
};
