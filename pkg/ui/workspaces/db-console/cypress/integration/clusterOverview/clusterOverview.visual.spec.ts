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

describe("Cluster Overview - check the number of live nodes and total ranges", () => {
  it("renders default view", () => {
    cy.visit("#/overview");
    cy.findAllByText("Capacity Usage");
    cy.findByText("Node Status");
    cy.findByText("Replication Status");
    cy.log("find the number of live nodes");
    cy.findAllByText("1");
    cy.log("find the number of ranges in replication status");
    cy.findAllByText("39");
  });
});

describe("Cluster Overview - Check Node map", () => {
  it("checks node maps information", () => {
    cy.visit("#/overview");
    cy.findAllByText("Node List").eq(1).click({ force: true });
    cy.findAllByText("Node Map").click();
    cy.log("Show the Node Map");
    cy.findByText("View the Node Map");
  });
});
