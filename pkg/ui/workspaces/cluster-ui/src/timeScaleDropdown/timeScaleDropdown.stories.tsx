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

import React, { useState } from "react";
import { storiesOf } from "@storybook/react";
import { TimeScaleDropdown } from "./timeScaleDropdown";
import { defaultTimeScaleOptions, defaultTimeScaleSelected } from "./utils";
import moment from "moment";

export function TimeScaleDropdownWrapper({
  initialTimeScale = defaultTimeScaleSelected,
}): React.ReactElement {
  const [timeScale, setTimeScale] = useState(initialTimeScale);
  return (
    <TimeScaleDropdown currentScale={timeScale} setTimeScale={setTimeScale} />
  );
}

storiesOf("TimeScaleDropdown", module)
  .add("default", () => <TimeScaleDropdownWrapper />)
  .add("custom", () => (
    <TimeScaleDropdownWrapper
      initialTimeScale={{
        sampleSize: defaultTimeScaleOptions["Past 6 Hours"].sampleSize,
        windowSize: moment.duration(6, "h"),
        fixedWindowEnd: moment().subtract(10, "m"),
        key: "Custom",
      }}
    />
  ));
