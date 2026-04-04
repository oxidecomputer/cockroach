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
import { randomName } from "src/storybook/fixtures";
import { DatabasesPage, DatabasesPageProps } from "./databasesPage";

import * as H from "history";
const history = H.createHashHistory();

const withLoadingIndicator: DatabasesPageProps = {
  loading: true,
  loaded: false,
  lastError: undefined,
  automaticStatsCollectionEnabled: true,
  databases: [],
  sortSetting: {
    ascending: false,
    columnTitle: "name",
  },
  onSortingChange: () => {},
  refreshDatabases: () => {},
  refreshSettings: () => {},
  refreshDatabaseDetails: () => {},
  refreshTableStats: () => {},
  location: history.location,
  history,
  match: {
    url: "",
    path: history.location.pathname,
    isExact: false,
    params: {},
  },
};

const withoutData: DatabasesPageProps = {
  loading: false,
  loaded: true,
  lastError: null,
  automaticStatsCollectionEnabled: true,
  databases: [],
  sortSetting: {
    ascending: false,
    columnTitle: "name",
  },
  onSortingChange: () => {},
  refreshDatabases: () => {},
  refreshSettings: () => {},
  refreshDatabaseDetails: () => {},
  refreshTableStats: () => {},
  location: history.location,
  history,
  match: {
    url: "",
    path: history.location.pathname,
    isExact: false,
    params: {},
  },
};

const withData: DatabasesPageProps = {
  loading: false,
  loaded: true,
  lastError: null,
  showNodeRegionsColumn: true,
  automaticStatsCollectionEnabled: true,
  sortSetting: {
    ascending: false,
    columnTitle: "name",
  },
  databases: _.map(Array(42), _item => {
    return {
      loading: false,
      loaded: true,
      lastError: null,
      name: randomName(),
      sizeInBytes: _.random(1000.0) * 1024 ** _.random(1, 2),
      tableCount: _.random(5, 100),
      rangeCount: _.random(50, 500),
      missingTables: [],
      nodesByRegionString:
        "gcp-europe-west1(n8), gcp-us-east1(n1), gcp-us-west1(n6)",
    };
  }),
  onSortingChange: () => {},
  refreshDatabases: () => {},
  refreshSettings: () => {},
  refreshDatabaseDetails: () => {},
  refreshTableStats: () => {},
  location: history.location,
  history,
  match: {
    url: "",
    path: history.location.pathname,
    isExact: false,
    params: {},
  },
};

storiesOf("Databases Page", module)
  .addDecorator(withRouterProvider)
  .addDecorator(withBackground)
  .add("with data", () => <DatabasesPage {...withData} />)
  .add("without data", () => <DatabasesPage {...withoutData} />)
  .add("with loading indicator", () => (
    <DatabasesPage {...withLoadingIndicator} />
  ));
