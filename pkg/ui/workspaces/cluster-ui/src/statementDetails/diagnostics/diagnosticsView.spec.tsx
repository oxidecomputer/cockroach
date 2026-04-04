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
import { mount, ReactWrapper } from "enzyme";
import sinon from "sinon";
import Long from "long";
import { MemoryRouter } from "react-router-dom";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { Button } from "@cockroachlabs/ui-components";

import { DiagnosticsView } from "./diagnosticsView";
import { Table } from "src/table";
import { TestStoreProvider } from "src/test-utils";

type IStatementDiagnosticsReport = cockroach.server.serverpb.IStatementDiagnosticsReport;

const sandbox = sinon.createSandbox();

const activateDiagnosticsRef = { current: { showModalFor: jest.fn() } };

function generateDiagnosticsRequest(
  extendObject: Partial<IStatementDiagnosticsReport> = {},
): IStatementDiagnosticsReport {
  const diagnosticsRequest = {
    statement_fingerprint: "SELECT * FROM table",
    completed: true,
    requested_at: {
      seconds: Long.fromNumber(Date.now()),
      nanos: Math.random() * 1000000,
    },
  };
  Object.assign(diagnosticsRequest, extendObject);
  return diagnosticsRequest;
}

describe("DiagnosticsView", () => {
  let wrapper: ReactWrapper;
  const statementFingerprint = "some-id";

  beforeEach(() => {
    sandbox.reset();
  });

  describe("With Empty state", () => {
    beforeEach(() => {
      wrapper = mount(
        <MemoryRouter>
          <DiagnosticsView
            activateDiagnosticsRef={activateDiagnosticsRef}
            statementFingerprint={statementFingerprint}
            hasData={false}
            diagnosticsReports={[]}
            dismissAlertMessage={() => {}}
          />
        </MemoryRouter>,
      );
    });

    it("opens the statement diagnostics modal when Activate button is clicked", () => {
      const activateButtonComponent = wrapper.find(Button).first();
      activateButtonComponent.simulate("click");
      expect(activateDiagnosticsRef.current.showModalFor).toBeCalledWith(
        statementFingerprint,
      );
    });
  });

  describe("With tracing data", () => {
    beforeEach(() => {
      const diagnosticsRequests: IStatementDiagnosticsReport[] = [
        generateDiagnosticsRequest(),
        generateDiagnosticsRequest(),
      ];

      wrapper = mount(
        <TestStoreProvider>
          <DiagnosticsView
            activateDiagnosticsRef={activateDiagnosticsRef}
            statementFingerprint={statementFingerprint}
            hasData={true}
            diagnosticsReports={diagnosticsRequests}
            dismissAlertMessage={() => {}}
          />
        </TestStoreProvider>,
      );
    });

    it("renders Table component when diagnostics data is provided", () => {
      assert.isTrue(wrapper.find(Table).exists());
    });

    it("opens the statement diagnostics modal when Activate button is clicked", () => {
      const activateButtonComponent = wrapper
        .findWhere(n => n.prop("children") === "Activate diagnostics")
        .first();
      activateButtonComponent.simulate("click");
      expect(activateDiagnosticsRef.current.showModalFor).toBeCalledWith(
        statementFingerprint,
      );
    });

    it("Activate button is hidden if diagnostics is requested and waiting query", () => {
      const diagnosticsRequests: IStatementDiagnosticsReport[] = [
        generateDiagnosticsRequest({ completed: false }),
        generateDiagnosticsRequest(),
      ];
      wrapper = mount(
        <TestStoreProvider>
          <DiagnosticsView
            activateDiagnosticsRef={activateDiagnosticsRef}
            statementFingerprint={statementFingerprint}
            hasData={true}
            diagnosticsReports={diagnosticsRequests}
            dismissAlertMessage={() => {}}
          />
        </TestStoreProvider>,
      );
      const activateButtonComponent = wrapper
        .findWhere(n => n.prop("children") === "Activate diagnostics")
        .first();
      assert.isFalse(activateButtonComponent.exists());
    });

    it("Cancel request button shows if diagnostics is requested and waiting query", () => {
      const diagnosticsRequests: IStatementDiagnosticsReport[] = [
        generateDiagnosticsRequest({ completed: false }),
        generateDiagnosticsRequest(),
      ];
      wrapper = mount(
        <TestStoreProvider>
          <DiagnosticsView
            activateDiagnosticsRef={activateDiagnosticsRef}
            statementFingerprint={statementFingerprint}
            hasData={true}
            diagnosticsReports={diagnosticsRequests}
            dismissAlertMessage={() => {}}
          />
        </TestStoreProvider>,
      );
      const cancelButtonComponent = wrapper
        .findWhere(n => n.prop("children") === "Cancel request")
        .first();
      assert.isTrue(cancelButtonComponent.exists());
    });
  });
});
