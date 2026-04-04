// Copyright 2023 The Cockroach Authors.
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

export type UserSQLRolesRequestMessage = cockroach.server.serverpb.UserSQLRolesRequest;
export type UserSQLRolesResponseMessage = cockroach.server.serverpb.UserSQLRolesResponse;

export function getUserSQLRoles(): Promise<UserSQLRolesResponseMessage> {
  return fetchData(
    cockroach.server.serverpb.UserSQLRolesResponse,
    `/_status/sqlroles`,
    null,
    null,
    "30M",
  );
}
