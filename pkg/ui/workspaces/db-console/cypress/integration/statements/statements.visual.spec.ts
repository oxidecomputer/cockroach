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

describe("Statements - Check whether the table content has been rendered properly", () => {
  it("renders default view", () => {
    cy.visit("#/statements");
    cy.get("section").contains("Statements");
  });
});

describe("Statements - can successfully activate diagnostics", () => {
  it("Writes some queries and activates diagnostics", () => {
    cy.visit("#/statements");
    cy.exec('cockroach sql --insecure --execute="create table test (id int)";');
    cy.exec(
      'cockroach sql --insecure --execute="insert into test values (1)";'
    );
    cy.log("Searching for the specific and activating its diagnostics");
    cy.findAllByPlaceholderText("Search Statement")
      .should("exist")
      .click()
      .type("insert into test");
    cy.findAllByRole("button", "Enter").eq(0).click();
    cy.findAllByText("Statements");
    cy.log(
      "Check whether the queries contain 'Activate' button and then click"
    );
    cy.findByText("Activate").eq(0).click().should("exist");
    cy.log("Click the button for the diagnostics");
    cy.findAllByText("Activate").eq(1).click();
    cy.findAllByText("WAITING").should("exist");
    cy.exec(
      'cockroach sql --insecure --execute="insert into test values (1)";'
    );
    cy.reload();
    cy.findAllByText("Activate").should("exist");
  });
});

// TODO: Download and check the downloaded diagnostics file
