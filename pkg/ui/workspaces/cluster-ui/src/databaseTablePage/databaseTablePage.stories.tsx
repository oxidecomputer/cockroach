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

import React from "react";
import { storiesOf } from "@storybook/react";
import _ from "lodash";

import { withBackground, withRouterProvider } from "src/storybook/decorators";
import {
  randomName,
  randomRole,
  randomTablePrivilege,
} from "src/storybook/fixtures";
import { DatabaseTablePage, DatabaseTablePageProps } from "./databaseTablePage";
import moment from "moment";
import * as H from "history";
const history = H.createHashHistory();

const withLoadingIndicator: DatabaseTablePageProps = {
  databaseName: randomName(),
  name: randomName(),
  automaticStatsCollectionEnabled: true,
  details: {
    loading: true,
    loaded: false,
    lastError: undefined,
    createStatement: "",
    replicaCount: 0,
    indexNames: [],
    grants: [],
    statsLastUpdated: moment("0001-01-01T00:00:00Z"),
  },
  stats: {
    loading: true,
    loaded: false,
    lastError: undefined,
    sizeInBytes: 0,
    rangeCount: 0,
  },
  indexStats: {
    loading: true,
    loaded: false,
    lastError: undefined,
    stats: [],
    lastReset: moment("2021-09-04T13:55:00Z"),
  },
  location: history.location,
  history,
  match: {
    url: "",
    path: history.location.pathname,
    isExact: false,
    params: {},
  },
  refreshTableDetails: () => {},
  refreshTableStats: () => {},
  refreshIndexStats: () => {},
  resetIndexUsageStats: () => {},
  refreshSettings: () => {},
  refreshUserSQLRoles: () => {},
};

const name = randomName();

const withData: DatabaseTablePageProps = {
  databaseName: randomName(),
  name: name,
  automaticStatsCollectionEnabled: true,
  details: {
    loading: false,
    loaded: true,
    lastError: null,
    createStatement: `
      CREATE TABLE public.${name} (
        id UUID NOT NULL,
        city VARCHAR NOT NULL,
        name VARCHAR NULL,
        address VARCHAR NULL,
        credit_card VARCHAR NULL,
        CONSTRAINT "primary" PRIMARY KEY (city ASC, id ASC),
        FAMILY "primary" (id, city, name, address, credit_card)
      )
    `,
    replicaCount: 7,
    indexNames: _.map(Array(3), randomName),
    grants: _.uniq(
      _.map(Array(12), () => {
        return {
          user: randomRole(),
          privilege: randomTablePrivilege(),
        };
      }),
    ),
    statsLastUpdated: moment("0001-01-01T00:00:00Z"),
  },
  showNodeRegionsSection: true,
  stats: {
    loading: false,
    loaded: true,
    lastError: null,
    sizeInBytes: 44040192,
    rangeCount: 4200,
    nodesByRegionString:
      "gcp-europe-west1(n8), gcp-us-east1(n1), gcp-us-west1(n6)",
  },
  indexStats: {
    loading: false,
    loaded: true,
    lastError: null,
    stats: [
      {
        totalReads: 0,
        lastUsed: moment("2021-10-11T11:29:00Z"),
        lastUsedType: "read",
        indexName: "primary",
      },
      {
        totalReads: 3,
        lastUsed: moment("2021-11-10T16:29:00Z"),
        lastUsedType: "read",
        indexName: "primary",
      },
      {
        totalReads: 2,
        lastUsed: moment("2021-09-04T13:55:00Z"),
        lastUsedType: "reset",
        indexName: "secondary",
      },
    ],
    lastReset: moment("2021-09-04T13:55:00Z"),
  },
  location: history.location,
  history,
  match: {
    url: "",
    path: history.location.pathname,
    isExact: false,
    params: {},
  },
  refreshTableDetails: () => {},
  refreshTableStats: () => {},
  refreshIndexStats: () => {},
  resetIndexUsageStats: () => {},
  refreshSettings: () => {},
  refreshUserSQLRoles: () => {},
};

storiesOf("Database Table Page", module)
  .addDecorator(withRouterProvider)
  .addDecorator(withBackground)
  .add("with data", () => <DatabaseTablePage {...withData} />)
  .add("with loading indicator", () => (
    <DatabaseTablePage {...withLoadingIndicator} />
  ));
