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
import { assert } from "chai";
import { ReactWrapper, mount } from "enzyme";
import { MemoryRouter } from "react-router-dom";

import {
  filterBySearchQuery,
  StatementsPage,
  StatementsPageProps,
  StatementsPageState,
} from "src/statementsPage";
import statementsPagePropsFixture from "./statementsPage.fixture";
import { AggregateStatistics } from "../statementsTable";
import { FlatPlanNode } from "../statementDetails";

describe("StatementsPage", () => {
  describe("Statements table", () => {
    it("sorts data by Execution Count DESC as default option", () => {
      const rootWrapper = mount(
        <MemoryRouter>
          <StatementsPage {...statementsPagePropsFixture} />
        </MemoryRouter>,
      );

      const statementsPageWrapper: ReactWrapper<
        StatementsPageProps,
        StatementsPageState,
        React.Component<any, any>
      > = rootWrapper.find(StatementsPage).first();
      const statementsPageInstance = statementsPageWrapper.instance();

      assert.equal(
        statementsPageInstance.props.sortSetting.columnTitle,
        "executionCount",
      );
      assert.equal(statementsPageInstance.props.sortSetting.ascending, false);
    });
  });

  describe("filterBySearchQuery", () => {
    const testPlanNode: FlatPlanNode = {
      name: "render",
      attrs: [],
      children: [
        {
          name: "group (scalar)",
          attrs: [],
          children: [
            {
              name: "filter",
              attrs: [
                {
                  key: "filter",
                  values: ["variable = _"],
                  warn: false,
                },
              ],
              children: [
                {
                  name: "virtual table",
                  attrs: [
                    {
                      key: "table",
                      values: ["cluster_settings@primary"],
                      warn: false,
                    },
                  ],
                  children: [],
                },
              ],
            },
          ],
        },
      ],
    };

    const statement: AggregateStatistics = {
      aggregatedFingerprintID: "",
      aggregatedFingerprintHexID: "",
      aggregatedTs: 0,
      aggregationInterval: 0,
      database: "",
      applicationName: "",
      fullScan: false,
      implicitTxn: false,
      summary: "",
      label:
        "SELECT count(*) > _ FROM [SHOW ALL CLUSTER SETTINGS] AS _ (v) WHERE v = '_'",
      stats: {
        sensitive_info: {
          most_recent_plan_description: testPlanNode,
        },
      },
    };

    assert.equal(filterBySearchQuery(statement, "select"), true);
    assert.equal(filterBySearchQuery(statement, "count"), true);
    assert.equal(filterBySearchQuery(statement, "select count"), true);
    assert.equal(filterBySearchQuery(statement, "cluster settings"), true);

    // Searching by plan should be false.
    assert.equal(filterBySearchQuery(statement, "virtual table"), false);
    assert.equal(filterBySearchQuery(statement, "group (scalar)"), false);
    assert.equal(filterBySearchQuery(statement, "node_build_info"), false);
    assert.equal(filterBySearchQuery(statement, "crdb_internal"), false);
  });
});
