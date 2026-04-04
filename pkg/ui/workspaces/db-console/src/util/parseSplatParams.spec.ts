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

import { assert } from "chai";
import { createMemoryHistory, History } from "history";
import { match as Match } from "react-router-dom";
import { parseSplatParams } from "./parseSplatParams";

describe("parseSplatParams", () => {
  let history: History;
  let match: Match;

  beforeEach(() => {
    history = createMemoryHistory({ initialEntries: ["/"] });
    match = {
      path: "/",
      params: {},
      url: "http://localhost/",
      isExact: true,
    };
  });

  it("returns remaining part of location path", () => {
    history.push("/overview/map/region=us-west/zone=a");
    match.path = "/overview/map/";

    assert.equal(
      parseSplatParams(match, history.location),
      "region=us-west/zone=a",
    );
  });

  it("trims out leading / from remaining path", () => {
    history.push("/overview/map/region=us-west/zone=a");
    match.path = "/overview/map";

    assert.equal(
      parseSplatParams(match, history.location),
      "region=us-west/zone=a",
    );
  });

  it("returns empty string if path is fully matched", () => {
    history.push("/overview/map");
    match.path = "/overview/map";

    assert.equal(parseSplatParams(match, history.location), "");
  });
});
