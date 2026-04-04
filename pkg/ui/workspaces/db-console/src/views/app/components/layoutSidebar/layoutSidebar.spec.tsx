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

import React from "react";
import { shallow } from "enzyme";
import { createMemoryHistory, History } from "history";
import { match as Match } from "react-router";
import { assert } from "chai";
import "src/enzymeInit";
import { Sidebar } from "./index";

describe("LayoutSidebar", () => {
  let history: History;
  let match: Match;

  beforeEach(() => {
    history = createMemoryHistory();
    match = {
      isExact: true,
      params: {},
      path: "/reports/network",
      url: "",
    };
  });

  it("does not show Network Latency link for single node cluster", () => {
    const wrapper = shallow(
      <Sidebar
        history={history}
        match={match}
        location={history.location}
        isSingleNodeCluster={true}
      />,
    );
    assert.isFalse(
      wrapper.findWhere(w => w.prop("to") === "/reports/network").exists(),
    );
  });

  it("shows Network Latency link for multi node cluster", () => {
    const wrapper = shallow(
      <Sidebar
        history={history}
        match={match}
        location={history.location}
        isSingleNodeCluster={false}
      />,
    );
    assert.isTrue(
      wrapper.findWhere(w => w.prop("to") === "/reports/network").exists(),
    );
  });
});
