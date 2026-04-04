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

import { EmptyTable } from "./emptyTable";
import emptyListResultsImg from "../../assets/emptyState/empty-list-results.svg";
import notFoundImg from "../../assets/emptyState/not-found-404.svg";
import { Button } from "src/button";
import SpinIcon from "../../icon/spin";

storiesOf("EmptyTablePlaceholder", module)
  .add("default", () => <EmptyTable />)
  .add("with icon", () => <EmptyTable icon={emptyListResultsImg} />)
  .add("with message", () => (
    <EmptyTable
      icon={emptyListResultsImg}
      message={"Lorem ipsum dolor sit amet, consectetur adipiscing elit."}
    />
  ))
  .add("with button", () => (
    <EmptyTable
      icon={emptyListResultsImg}
      message={"Lorem ipsum dolor sit amet, consectetur adipiscing elit."}
      footer={<Button type="primary">Activate</Button>}
    />
  ))
  .add("no backup found", () => (
    <EmptyTable
      icon={notFoundImg}
      title="No backup found"
      message="The backup you’re looking for does not exist."
    />
  ))
  .add("with SVG element as icon", () => {
    return (
      <EmptyTable
        icon={<SpinIcon width={72} height={72} />}
        title="It's SVG icon"
      />
    );
  });
