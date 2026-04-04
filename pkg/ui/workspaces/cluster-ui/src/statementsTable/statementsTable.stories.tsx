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
import { MemoryRouter } from "react-router-dom";
import {
  makeStatementsColumns,
  StatementsSortedTable,
} from "./statementsTable";
import statementsPagePropsFixture from "src/statementsPage/statementsPage.fixture";
import { calculateTotalWorkload } from "src/util";

const { statements } = statementsPagePropsFixture;

storiesOf("StatementsSortedTable", module)
  .addDecorator(storyFn => <MemoryRouter>{storyFn()}</MemoryRouter>)
  .add("with data", () => (
    <StatementsSortedTable
      className="statements-table"
      data={statements}
      columns={makeStatementsColumns(
        statements,
        ["$ internal"],
        calculateTotalWorkload(statements),
        { "1": "gcp-europe-west1", "2": "gcp-us-east1", "3": "gcp-us-west1" },
        "statement",
        false,
        false,
        null,
        React.createRef(),
      )}
      sortSetting={{
        ascending: false,
        columnTitle: "rowsRead",
      }}
      pagination={{
        pageSize: 20,
        current: 1,
      }}
    />
  ))
  .add("with data and VIEWACTIVITYREDACTED role", () => (
    <StatementsSortedTable
      className="statements-table"
      data={statements}
      columns={makeStatementsColumns(
        statements,
        ["$ internal"],
        calculateTotalWorkload(statements),
        { "1": "gcp-europe-west1", "2": "gcp-us-east1", "3": "gcp-us-west1" },
        "statement",
        false,
        true,
        null,
        React.createRef(),
      )}
      sortSetting={{
        ascending: false,
        columnTitle: "rowsRead",
      }}
      pagination={{
        pageSize: 20,
        current: 1,
      }}
    />
  ))
  .add("empty table", () => <StatementsSortedTable empty />);
