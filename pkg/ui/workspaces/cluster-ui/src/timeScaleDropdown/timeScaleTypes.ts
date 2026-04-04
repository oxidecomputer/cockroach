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

import moment from "moment";
import { defaultTimeScaleOptions } from "./utils";

/**
 * TimeWindow represents an absolute window of time, defined with a start and
 * end time.
 */
export interface TimeWindow {
  start: moment.Moment;
  end: moment.Moment;
}

/**
 * TimeScale describes the requested dimensions, from which one can derive concrete TimeWindows using toDateRange.
 */
export interface TimeScale {
  /**
   * The key used to index in to the defaultTimeScaleOptions collection.
   * The key is "Custom" when specifying a custom time that is not one of the default options
   */
  key?: string;
  // The size of a global time window. Default is ten minutes.
  windowSize: moment.Duration;
  // The length of time the global time window is valid. Default is ten seconds.
  windowValid?: moment.Duration;
  // The expected duration of individual samples for queries at this time scale.
  sampleSize: moment.Duration;
  /**
   * The fixed end time of the window, or false if it should be a dynamically moving "now".
   * Typically, when the `key` property is a default option, `fixedWindowEnd` should be false.
   * And when the `key` property is "Custom" `fixedWindowEnd` should be a specific Moment.
   * It is unclear if there are legitimate reasons for the two being out of sync.
   */
  fixedWindowEnd: moment.Moment | false;
}

export class TimeScaleState {
  // Currently selected scale.
  scale: TimeScale;
  constructor() {
    this.scale = {
      ...defaultTimeScaleOptions["Past 10 Minutes"],
      fixedWindowEnd: false,
      key: "Past 10 Minutes",
    };
  }
}

export type TimeScaleOption = Omit<TimeScale, "fixedWindowEnd">;

export interface TimeScaleOptions {
  [key: string]: TimeScaleOption;
}

export enum ArrowDirection {
  LEFT,
  RIGHT,
  CENTER,
}
