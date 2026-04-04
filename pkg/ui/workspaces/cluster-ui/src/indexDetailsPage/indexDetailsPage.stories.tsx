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

import { withBackground, withRouterProvider } from "src/storybook/decorators";
import { randomName } from "src/storybook/fixtures";
import { IndexDetailsPage, IndexDetailsPageProps } from "./indexDetailsPage";
import moment from "moment";

const withData: IndexDetailsPageProps = {
  databaseName: randomName(),
  tableName: randomName(),
  indexName: randomName(),
  details: {
    loading: false,
    loaded: true,
    createStatement: `
      CREATE UNIQUE INDEX "primary" ON system.public.database_role_settings USING btree (database_id ASC, role_name ASC)
    `,
    totalReads: 0,
    lastRead: moment("2021-10-21T22:00:00Z"),
    lastReset: moment("2021-12-02T07:12:00Z"),
  },
  refreshIndexStats: () => {},
  resetIndexUsageStats: () => {},
  refreshNodes: () => {},
  refreshUserSQLRoles: () => {},
};

storiesOf("Index Details Page", module)
  .addDecorator(withRouterProvider)
  .addDecorator(withBackground)
  .add("with data", () => <IndexDetailsPage {...withData} />);
