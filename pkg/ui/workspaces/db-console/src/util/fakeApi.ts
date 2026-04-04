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

import * as $protobuf from "protobufjs";

import { cockroach } from "src/js/protos";
import { API_PREFIX, STATUS_PREFIX, toArrayBuffer } from "src/util/api";
import fetchMock from "src/util/fetch-mock";

const {
  DatabasesResponse,
  DatabaseDetailsResponse,
  SettingsResponse,
  TableDetailsResponse,
  TableStatsResponse,
  TableIndexStatsResponse,
} = cockroach.server.serverpb;

// These test-time functions provide typesafe wrappers around fetchMock,
// stubbing HTTP responses from the admin API.
//
// Be sure to call `restore()` after each test that uses `fakeApi`,
// so that your stubbed responses won't leak out to other tests.
//
// Example usage:
//
//   describe("The thing I'm testing", function() {
//     it("works like this", function() {
//       // 1. Set up a fake response from the /databases endpoint.
//       fakeApi.stubDatabases({
//         databases: ["one", "two", "three"],
//       });
//
//       // 2. Run your code that hits the /databases endpoint.
//       // ...
//
//       // 3. Assert on its data / behavior.
//       assert.deepEqual(myThing.databaseNames(), ["one", "two", "three"]);
//     });
//
//     // 4. Tear down any fake responses we set up.
//     afterEach(function() {
//       fakeApi.restore();
//     });
//   });

export function restore() {
  fetchMock.restore();
}

export function stubClusterSettings(
  response: cockroach.server.serverpb.ISettingsResponse,
) {
  stubGet(
    "/settings?unredacted_values=true",
    SettingsResponse.encode(response),
    API_PREFIX,
  );
}

export function stubDatabases(
  response: cockroach.server.serverpb.IDatabasesResponse,
) {
  stubGet("/databases", DatabasesResponse.encode(response), API_PREFIX);
}

export function stubDatabaseDetails(
  database: string,
  response: cockroach.server.serverpb.IDatabaseDetailsResponse,
) {
  stubGet(
    `/databases/${database}?include_stats=true`,
    DatabaseDetailsResponse.encode(response),
    API_PREFIX,
  );
}

export function stubTableDetails(
  database: string,
  table: string,
  response: cockroach.server.serverpb.ITableDetailsResponse,
) {
  stubGet(
    `/databases/${database}/tables/${table}`,
    TableDetailsResponse.encode(response),
    API_PREFIX,
  );
}

export function stubTableStats(
  database: string,
  table: string,
  response: cockroach.server.serverpb.ITableStatsResponse,
) {
  stubGet(
    `/databases/${database}/tables/${table}/stats`,
    TableStatsResponse.encode(response),
    API_PREFIX,
  );
}

export function stubIndexStats(
  database: string,
  table: string,
  response: cockroach.server.serverpb.ITableIndexStatsResponse,
) {
  stubGet(
    `/databases/${database}/tables/${table}/indexstats`,
    TableIndexStatsResponse.encode(response),
    STATUS_PREFIX,
  );
}

function stubGet(path: string, writer: $protobuf.Writer, prefix: string) {
  fetchMock.get(`${prefix}${path}`, toArrayBuffer(writer.finish()));
}
