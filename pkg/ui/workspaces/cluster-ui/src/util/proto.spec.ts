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

import { INodeStatus, rollupStoreMetrics } from "./proto";

describe("Proto utils", () => {
  describe("rollupStoreMetrics", () => {
    let nodeStatus: Partial<INodeStatus>;
    let statusWithRolledMetrics: Partial<INodeStatus>;

    beforeEach(() => {
      nodeStatus = {
        metrics: {
          a: 10,
          b: 5,
          c: 0,
          y: 15,
          z: 5,
        },
        store_statuses: [
          {
            metrics: {
              c: 25,
              d: 5,
              e: 5,
            },
          },
          {
            metrics: {
              a: 5,
              b: 100,
              x: 0,
              y: 20,
              z: 0,
            },
          },
        ],
      };
      statusWithRolledMetrics = {
        ...nodeStatus,
        metrics: {
          a: 15,
          b: 105,
          c: 25,
          d: 5,
          e: 5,
          x: 0,
          y: 35,
          z: 5,
        },
      };
    });

    it("sums up values for every metric", () => {
      assert.deepEqual(
        rollupStoreMetrics(nodeStatus),
        statusWithRolledMetrics.metrics,
      );
    });
  });
});
