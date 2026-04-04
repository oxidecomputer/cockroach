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
import { MemoryRouter } from "react-router-dom";
import { cloneDeep } from "lodash";

import { StatementsPage } from "./statementsPage";
import statementsPagePropsFixture, {
  statementsPagePropsWithRequestError,
} from "./statementsPage.fixture";

storiesOf("StatementsPage", module)
  .addDecorator(storyFn => <MemoryRouter>{storyFn()}</MemoryRouter>)
  .addDecorator(storyFn => (
    <div style={{ backgroundColor: "#F5F7FA" }}>{storyFn()}</div>
  ))
  .add("with data", () => <StatementsPage {...statementsPagePropsFixture} />)
  .add("without data", () => (
    <StatementsPage {...statementsPagePropsFixture} statements={[]} />
  ))
  .add("with empty search result", () => {
    const props = cloneDeep(statementsPagePropsFixture);
    const { history } = props;
    const searchParams = new URLSearchParams(history.location.search);
    searchParams.set("q", "aaaaaaa");
    history.location.search = searchParams.toString();
    return (
      <StatementsPage
        {...props}
        {...statementsPagePropsFixture}
        statements={[]}
        history={history}
      />
    );
  })
  .add("with error", () => {
    return (
      <StatementsPage
        {...statementsPagePropsWithRequestError}
        statements={[]}
      />
    );
  })
  .add("with loading state", () => {
    return <StatementsPage {...statementsPagePropsFixture} statements={null} />;
  });
