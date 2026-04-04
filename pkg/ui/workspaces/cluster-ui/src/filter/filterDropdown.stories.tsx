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

import React from "react";
import { storiesOf } from "@storybook/react";
import { noop } from "lodash";
import { FilterDropdown } from "./filterDropdown";
import { FilterCheckboxOption } from "./filterCheckboxOption";
import { FilterSearchOption } from "./filterSearchOption";

storiesOf("FilterDropdown", module)
  .addDecorator(renderChild => (
    <div style={{ padding: "12px", display: "flex" }}>{renderChild()}</div>
  ))
  .add("default", () => (
    <FilterDropdown label="Filters" onSubmit={noop}>
      <FilterCheckboxOption
        label="Node ID"
        value={[{ label: "1", value: "1" }]}
        options={[
          { label: "1", value: "1" },
          { label: "2", value: "2" },
          { label: "3", value: "3" },
          { label: "4", value: "4" },
        ]}
        placeholder="Select"
      />
      <FilterSearchOption label="Store ID" />
    </FilterDropdown>
  ));
