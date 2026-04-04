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

describe("Databases - show databases without any tables", () => {
  it("renders default view", () => {
    cy.visit("#/databases/tables");
    cy.findAllByText("Databases").should("exist");
    cy.findAllByText("defaultdb").should("exist");
  });
});

describe("Databases - show databases with tables", () => {
  it("renders the view for databases/tables view", () => {
    cy.visit("#/databases/tables");
    cy.exec(
      'cockroach sql --insecure --execute="create table if not exists test (id int)";'
    );
    cy.exec(
      'cockroach sql --insecure --execute="create table if not exists test1 (id int)";'
    );
    cy.exec(
      'cockroach sql --insecure --execute="create table if not exists test2 (id int)";'
    );
    cy.findAllByText("Load stats for all tables").should("exist").eq(0).click();
    cy.reload();
    cy.log("check whether the tables exit in the database");
    cy.findAllByText("public.test").should("exist");
    cy.findAllByText("public.test1").should("exist");
  });
});
