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
import { noop } from "lodash";

import { Dropdown, DropdownOption } from "./dropdown";
import { Button } from "src/button";
import { Download } from "@cockroachlabs/icons";

const items: DropdownOption[] = [
  { name: "A", value: "a" },
  { name: "B", value: "b" },
  { name: "C", value: "c" },
];

storiesOf("Dropdown", module)
  .addDecorator(renderChild => (
    <div style={{ padding: "12px", display: "flex" }}>{renderChild()}</div>
  ))
  .add("default", () => (
    <Dropdown onChange={noop} items={items}>
      Select
    </Dropdown>
  ))
  .add("with custom toggle icon", () => (
    <Dropdown
      onChange={noop}
      items={items}
      customToggleButton={
        <Button type="primary" textAlign="center">
          <Download />
        </Button>
      }
    />
  ))
  .add("with custom toggle button options", () => (
    <Dropdown
      onChange={noop}
      items={items}
      customToggleButtonOptions={{
        iconPosition: "left",
        size: "small",
        type: "unstyled-link",
      }}
    >
      Select options
    </Dropdown>
  ));
