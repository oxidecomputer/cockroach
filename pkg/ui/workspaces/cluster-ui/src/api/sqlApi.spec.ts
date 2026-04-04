// Copyright 2022 The Cockroach Authors.
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

import { SqlExecutionResponse, sqlResultsAreEmpty } from "./sqlApi";

describe("sqlApi", () => {
  test("sqlResultsAreEmpty should return true when there are no rows in the response", () => {
    const testCases: {
      response: SqlExecutionResponse<unknown>;
      expected: boolean;
    }[] = [
      {
        response: {
          num_statements: 1,
          execution: {
            retries: 0,
            txn_results: [
              {
                statement: 1,
                tag: "SELECT",
                start: "start-date",
                end: "end-date",
                rows_affected: 0,
                rows: [{ hello: "world" }],
              },
            ],
          },
        },
        expected: false,
      },
      {
        response: {
          num_statements: 1,
          execution: {
            retries: 0,
            txn_results: [
              {
                statement: 1,
                tag: "SELECT",
                start: "start-date",
                end: "end-date",
                rows_affected: 0,
                rows: [],
              },
            ],
          },
        },
        expected: true,
      },
      {
        response: {
          num_statements: 1,
          execution: {
            retries: 0,
            txn_results: [
              {
                statement: 1,
                tag: "SELECT",
                start: "start-date",
                end: "end-date",
                rows_affected: 0,
                columns: [],
              },
            ],
          },
        },
        expected: true,
      },
      {
        response: {},
        expected: true,
      },
    ];

    testCases.forEach(tc => {
      expect(sqlResultsAreEmpty(tc.response)).toEqual(tc.expected);
    });
  });
});
