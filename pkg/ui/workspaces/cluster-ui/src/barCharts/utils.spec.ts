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

import { assert } from "chai";
import { normalizeClosedDomain } from "./utils";

describe("barCharts utils", () => {
  describe("normalizeClosedDomain", () => {
    it("returns input args if domain values are not equal", () => {
      assert.deepStrictEqual(normalizeClosedDomain([10, 15]), [10, 15]);
    });

    it("returns increased end range by 1 if input start and end values are equal", () => {
      assert.deepStrictEqual(normalizeClosedDomain([10, 10]), [10, 11]);
    });
  });
});
