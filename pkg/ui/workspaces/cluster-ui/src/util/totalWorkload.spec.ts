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
import { calculateTotalWorkload } from "./totalWorkload";
import { aggStatFix } from "./totalWorkload.fixture";

describe("Calculating total workload", () => {
  it("calculating total workload with one statement", () => {
    const result = calculateTotalWorkload([aggStatFix]);
    // Using approximately because float handling by javascript is imprecise
    assert.approximately(result, 48.421019, 0.0000001);
  });

  it("calculating total workload with no statements", () => {
    const result = calculateTotalWorkload([]);
    assert.equal(result, 0);
  });

  it("calculating total workload with multiple statements", () => {
    const result = calculateTotalWorkload([aggStatFix, aggStatFix, aggStatFix]);
    // Using approximately because float handling by javascript is imprecise
    assert.approximately(result, 145.263057, 0.0000001);
  });
});
