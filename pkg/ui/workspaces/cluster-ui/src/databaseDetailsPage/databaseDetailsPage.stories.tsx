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
import {
  DatabaseDetailsPage,
  DatabaseDetailsPageProps,
  ViewMode,
} from "./databaseDetailsPage";

import * as H from "history";
import moment from "moment";
const history = H.createHashHistory();

const withLoadingIndicator: DatabaseDetailsPageProps = {
  loading: true,
  loaded: false,
  lastError: undefined,
  name: randomName(),
  tables: [],
  viewMode: ViewMode.Tables,
  sortSettingTables: {
    ascending: false,
    columnTitle: "name",
  },
  sortSettingGrants: {
    ascending: false,
    columnTitle: "name",
  },
  onSortingTablesChange: () => {},
  onSortingGrantsChange: () => {},
  refreshDatabaseDetails: () => {},
  refreshTableDetails: () => {},
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

const withoutData: DatabaseDetailsPageProps = {
  loading: false,
  loaded: true,
  lastError: null,
  name: randomName(),
  tables: [],
  viewMode: ViewMode.Tables,
  sortSettingTables: {
    ascending: false,
    columnTitle: "name",
  },
  sortSettingGrants: {
    ascending: false,
    columnTitle: "name",
  },
  onSortingTablesChange: () => {},
  onSortingGrantsChange: () => {},
  refreshDatabaseDetails: () => {},
  refreshTableDetails: () => {},
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

const withData: DatabaseDetailsPageProps = {
  loading: false,
  loaded: true,
  lastError: null,
  name: randomName(),
  tables: _.map(Array(42), _item => {
    const roles = _.uniq(
      _.map(new Array(_.random(1, 3)), _item => randomRole()),
    );
    const grants = _.uniq(
      _.map(new Array(_.random(1, 5)), _item => randomTablePrivilege()),
    );

    return {
      name: randomName(),
      details: {
        loading: false,
        loaded: true,
        lastError: null,
        columnCount: _.random(5, 42),
        indexCount: _.random(1, 6),
        userCount: roles.length,
        roles: roles,
        grants: grants,
        statsLastUpdated: moment("0001-01-01T00:00:00Z"),
      },
      showNodeRegionsColumn: true,
      stats: {
        loading: false,
        loaded: true,
        lastError: null,
        replicationSizeInBytes: _.random(1000.0) * 1024 ** _.random(1, 2),
        rangeCount: _.random(50, 500),
        nodesByRegionString:
          "gcp-europe-west1(n8), gcp-us-east1(n1), gcp-us-west1(n6)",
      },
    };
  }),
  viewMode: ViewMode.Tables,
  sortSettingTables: {
    ascending: false,
    columnTitle: "name",
  },
  sortSettingGrants: {
    ascending: false,
    columnTitle: "name",
  },
  onSortingTablesChange: () => {},
  onSortingGrantsChange: () => {},
  refreshDatabaseDetails: () => {},
  refreshTableDetails: () => {},
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

storiesOf("Database Details Page", module)
  .addDecorator(withRouterProvider)
  .addDecorator(withBackground)
  .add("with data", () => <DatabaseDetailsPage {...withData} />)
  .add("without data", () => <DatabaseDetailsPage {...withoutData} />)
  .add("with loading indicator", () => (
    <DatabaseDetailsPage {...withLoadingIndicator} />
  ));
