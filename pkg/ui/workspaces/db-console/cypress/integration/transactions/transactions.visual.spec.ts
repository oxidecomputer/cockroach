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

describe("Transactions - Check whether the table content has been rendered properly", () => {
  it("renders default view", () => {
    cy.visit("#/transactions");
    cy.findAllByText("Transactions").should("exist");
  });
});

describe("Transactions - Check the transactions statistics show up properly", () => {
  it("Writes some queries and check transactions", () => {
    cy.visit("#/transactions");
    cy.exec(
      'cockroach sql --insecure --execute="BEGIN; CREATE TABLE if not exists cypress_test (id int); INSERT INTO cypress_test VALUES (234); INSERT INTO cypress_test VALUES (234); SELECT * FROM cypress_test; COMMIT;"'
    );
    cy.reload();
    cy.findByPlaceholderText("Search transactions")
      .should("exist")
      .click()
      .type("CREATE TABLE cypress_test");
    cy.findAllByRole("button", "Enter").click();
    cy.findAllByText("cypress_test").should("exist").click();
    cy.log(
      "Click into the particular transaction and check if it has been executed properly"
    );
    cy.findByText("Transaction Details").should("exist");
    cy.findAllByText("CREATE TABLE IF").should("exist");
    cy.findAllByText("INSERT INTO cypress_test").should("exist");
    cy.findAllByText("SELECT FROM cypress_test").should("exist");
  });
});
