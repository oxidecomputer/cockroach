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
import { Moment } from "moment";
import { Text, TextTypes } from "src/components";

export interface DateRangeLabelProps {
  from: Moment;
  to: Moment;
}

export const DateRangeLabel: React.FC<DateRangeLabelProps> = ({ from, to }) => {
  const dateFormat = "MMM D";
  const timeFormat = "LT";
  const fromDateStr = from.format(dateFormat);
  const toDateStr = to.format(dateFormat);
  const fromTimeStr = from.format(timeFormat);
  const toTimeStr = to.format(timeFormat);
  const isUTC = to.isUTC() && from.isUTC();
  return (
    <div style={{ textAlign: "left" }}>
      <Text textType={TextTypes.Body}>
        {fromDateStr}
        {", "}
      </Text>
      <Text textType={TextTypes.BodyStrong}>{fromTimeStr}</Text>
      <Text textType={TextTypes.Body}>{" — "}</Text>
      <Text textType={TextTypes.Body}>
        {toDateStr}
        {", "}
      </Text>
      <Text textType={TextTypes.BodyStrong}>{toTimeStr}</Text>
      {isUTC && <Text textType={TextTypes.Body}>{" UTC"}</Text>}
    </div>
  );
};
