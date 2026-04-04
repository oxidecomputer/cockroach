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
import { storiesOf } from "@storybook/react";

import { InlineAlert } from "./inlineAlert";
import { styledWrapper } from "src/util/decorators";
import { Anchor } from "src/components";

storiesOf("InlineAlert", module)
  .addDecorator(styledWrapper({ padding: "24px" }))
  .add("with text title", () => (
    <InlineAlert title="Hello world!" message="blah-blah-blah" />
  ))
  .add("with Error intent", () => (
    <InlineAlert title="Hello world!" message="blah-blah-blah" intent="error" />
  ))
  .add("with link in title", () => (
    <InlineAlert
      title={
        <span>
          You do not have permission to view this information.{" "}
          <Anchor href="#">Learn more.</Anchor>
        </span>
      }
    />
  ))
  .add("with multiline message", () => (
    <InlineAlert
      title="Hello world!"
      message={
        <div>
          <div>Message 1</div>
          <div>Message 2</div>
          <div>Message 3</div>
        </div>
      }
    />
  ));
