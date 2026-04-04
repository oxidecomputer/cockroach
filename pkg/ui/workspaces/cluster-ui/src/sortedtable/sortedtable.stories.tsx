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
import { storiesOf } from "@storybook/react";

import { SortedTable } from "./";

storiesOf("Sorted table", module)
  .add("Empty state", () => <SortedTable empty />)
  .add("With data", () => {
    const columns = [
      {
        name: "Col 1",
        title: "Col 1",
        cell: (idx: number) => `row-${idx} col-1`,
      },
      {
        name: "Col 2",
        title: "Col 2",
        cell: (idx: number) => `row-${idx} col-2`,
      },
      {
        name: "Col 3",
        title: "Col 3",
        cell: (idx: number) => `row-${idx} col-3`,
      },
    ];
    return <SortedTable columns={columns} data={[1, 2, 3]} />;
  });
