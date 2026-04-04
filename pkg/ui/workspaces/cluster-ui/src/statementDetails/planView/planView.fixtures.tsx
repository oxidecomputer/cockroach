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

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { GlobalPropertiesType } from "./planView";

type IExplainTreePlanNode = cockroach.sql.IExplainTreePlanNode;

export const globalProperties: GlobalPropertiesType = {
  distribution: true,
  vectorized: true,
};

export const logicalPlan: IExplainTreePlanNode = {
  name: "root",
  attrs: [],
  children: [
    {
      name: "count",
      attrs: [],
      children: [
        {
          name: "upsert",
          attrs: [
            {
              key: "into",
              value:
                "vehicle_location_histories(city, ride_id, timestamp, lat, long)",
            },
            {
              key: "strategy",
              value: "opt upserter",
            },
          ],
          children: [
            {
              name: "buffer node",
              attrs: [
                {
                  key: "label",
                  value: "buffer 1",
                },
              ],
              children: [
                {
                  name: "row source to plan node",
                  attrs: [],
                  children: [
                    {
                      name: "render",
                      attrs: [
                        {
                          key: "render",
                          value: "column1",
                        },
                        {
                          key: "render",
                          value: "column2",
                        },
                        {
                          key: "render",
                          value: "column3",
                        },
                        {
                          key: "render",
                          value: "column4",
                        },
                        {
                          key: "render",
                          value: "column5",
                        },
                        {
                          key: "render",
                          value: "column4",
                        },
                        {
                          key: "render",
                          value: "column5",
                        },
                      ],
                      children: [
                        {
                          name: "values",
                          attrs: [
                            {
                              key: "size",
                              value: "5 columns, 1 row",
                            },
                            {
                              key: "row 0, expr",
                              value: "_",
                            },
                            {
                              key: "row 0, expr",
                              value: "_",
                            },
                            {
                              key: "row 0, expr",
                              value: "now()",
                            },
                            {
                              key: "row 0, expr",
                              value: "_",
                            },
                            {
                              key: "row 0, expr",
                              value: "_",
                            },
                          ],
                          children: [],
                        },
                      ],
                    },
                  ],
                },
              ],
            },
          ],
        },
      ],
    },
    {
      name: "postquery",
      attrs: [],
      children: [
        {
          name: "error if rows",
          attrs: [],
          children: [
            {
              name: "row source to plan node",
              attrs: [],
              children: [
                {
                  name: "lookup-join",
                  attrs: [
                    {
                      key: "table",
                      value: "rides@primary",
                    },
                    {
                      key: "type",
                      value: "anti",
                    },
                    {
                      key: "equality",
                      value: "(column1, column2) = (city, id)",
                    },
                    {
                      key: "equality cols are key",
                      value: "",
                    },
                    {
                      key: "parallel",
                      value: "",
                    },
                  ],
                  children: [
                    {
                      name: "render",
                      attrs: [
                        {
                          key: "render",
                          value: "column1",
                        },
                        {
                          key: "render",
                          value: "column2",
                        },
                      ],
                      children: [
                        {
                          name: "scan buffer node",
                          children: [],
                          attrs: [
                            {
                              key: "label",
                              value: "buffer 1",
                            },
                          ],
                        },
                      ],
                    },
                  ],
                },
              ],
            },
          ],
        },
      ],
    },
  ],
};
