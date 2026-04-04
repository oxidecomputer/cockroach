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
import { withRouterDecorator } from "src/util/decorators";
import { LoginPage } from "./loginPage";
import {
  loginPagePropsFixture,
  loginPagePropsLoadingFixture,
  loginPagePropsErrorFixture,
} from "./loginPage.fixture";

storiesOf("LoginPage", module)
  .addDecorator(withRouterDecorator)
  .add("Default state", () => <LoginPage {...loginPagePropsFixture} />)
  .add("Loading state", () => <LoginPage {...loginPagePropsLoadingFixture} />)
  .add("Error state", () => <LoginPage {...loginPagePropsErrorFixture} />);
