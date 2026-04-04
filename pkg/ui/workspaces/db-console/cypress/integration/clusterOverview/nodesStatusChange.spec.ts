// Copyright 2020 The Cockroach Authors.
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

import { TIME_UNTIL_NODE_DEAD } from "../../support/constants";

describe("Nodes status change", () => {
  beforeEach(() => {
    cy.teardown();
    cy.startCluster();
    cy.visit("#/overview/list");
  });

  describe("from Live to Decommissioned", () => {
    it("changes to Decommissioned status", () => {
      const nodeIdToDecommission = 2;
      cy
        .log("Validate all 4 nodes are Live")
        .get(
          ".nodes-overview__live-nodes-table tbody tr td span span:contains(Live)",
          {logMessage: "Node List table > rows with Live status"},
        )
        .should("have.length", 4);

      cy
        .log("Cluster summary section displays 4 live nodes")
        .get(
          ".node-liveness.cluster-summary__metric.live-nodes",
          {logMessage: "Cluster Summary > get Live Nodes value"},
        )
        .should("contain", 4);

      cy.decommissionNode(nodeIdToDecommission);
      cy.stopNode(nodeIdToDecommission);

      cy
        .log("Validate that only 3 live nodes remain")
        .get(
          ".nodes-overview__live-nodes-table tbody tr td span span:contains(Live)",
          { logMessage: "Node List table > rows with Live status" },
        )
        .should("have.length", 3);

      cy
        .log("...and 1 node is Decommissioning table")
        .get(
          ".nodes-overview__live-nodes-table",
          {logMessage: "Node List table > rows with Decommissioning status"},
        )
        .find("tbody tr:contains(Decommissioning)")
        .should("have.length", 1);

      cy
        .log("...and Cluster summary shows that 1 node is suspected")
        .get(
          ".node-liveness.cluster-summary__metric.suspect-nodes",
          {logMessage: "Cluster Summary > get Suspected Nodes value"},
        )
        .should("contain", 1);

      cy.wait(TIME_UNTIL_NODE_DEAD);

      cy
        .log("Validate that Decommissioned node table exists and contains one record")
        .get(
          ".nodes-overview__decommissioned-nodes-table",
          { logMessage: "Decommissioned Nodes table exists" },
        )
        .should("exist")
        .contains("tbody tr .status-column.status-column--color-decommissioned > span", "Decommissioned")
        .should("exist");
    });
  });
});
