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

import { configure } from "@storybook/react";

// Import global styles here
import "nvd3/build/nv.d3.min.css";
import "react-select/dist/react-select.css";
import "antd/es/tooltip/style/css";
import "styl/app.styl";
import "./styles.css";
import "src/views/app/containers/layout/layout.styl";

const req = require.context("../src/", true, /.stories.tsx$/);

function loadStories() {
  req.keys().forEach(filename => req(filename));
}

configure(loadStories, module);
