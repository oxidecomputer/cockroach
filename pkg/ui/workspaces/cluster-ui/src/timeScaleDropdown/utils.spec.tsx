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

import { TimeScale } from "./timeScaleTypes";
import moment from "moment";
import {
  defaultTimeScaleOptions,
  findClosestTimeScale,
  toDateRange,
  toRoundedDateRange,
} from "./utils";
import { assert } from "chai";

describe("timescale utils", (): void => {
  describe("toDateRange", () => {
    it("get date range", () => {
      const ts: TimeScale = {
        windowSize: moment.duration(5, "day"),
        sampleSize: moment.duration(5, "minutes"),
        fixedWindowEnd: moment.utc("2022.01.10 13:42"),
        key: "Custom",
      };
      const [start, end] = toDateRange(ts);
      assert.equal(start.format("YYYY.MM.DD HH:mm:ss"), "2022.01.05 13:42:00");
      assert.equal(end.format("YYYY.MM.DD HH:mm:ss"), "2022.01.10 13:42:00");
    });
  });

  describe("toRoundedDateRange", () => {
    it("round values", () => {
      const ts: TimeScale = {
        windowSize: moment.duration(5, "day"),
        sampleSize: moment.duration(5, "minutes"),
        fixedWindowEnd: moment.utc("2022.01.10 13:42"),
        key: "Custom",
      };
      const [start, end] = toRoundedDateRange(ts);
      assert.equal(start.format("YYYY.MM.DD HH:mm:ss"), "2022.01.05 13:00:00");
      assert.equal(end.format("YYYY.MM.DD HH:mm:ss"), "2022.01.10 13:59:59");
    });

    it("already rounded values", () => {
      const ts: TimeScale = {
        windowSize: moment.duration(5, "day"),
        sampleSize: moment.duration(5, "minutes"),
        fixedWindowEnd: moment.utc("2022.01.10 13:00"),
        key: "Custom",
      };
      const [start, end] = toRoundedDateRange(ts);
      assert.equal(start.format("YYYY.MM.DD HH:mm:ss"), "2022.01.05 13:00:00");
      assert.equal(end.format("YYYY.MM.DD HH:mm:ss"), "2022.01.10 13:59:59");
    });
  });

  describe("findClosestTimeScale", () => {
    it("should find the correct time scale", () => {
      // `seconds` != window size of any of the default options, `startSeconds` not specified.
      assert.deepEqual(findClosestTimeScale(defaultTimeScaleOptions, 15), {
        ...defaultTimeScaleOptions["Past 10 Minutes"],
        key: "Custom",
      });
      // `seconds` != window size of any of the default options, `startSeconds` not specified.
      assert.deepEqual(
        findClosestTimeScale(
          defaultTimeScaleOptions,
          moment.duration(moment().daysInMonth() * 5, "days").asSeconds(),
        ),
        { ...defaultTimeScaleOptions["Past 2 Months"], key: "Custom" },
      );
      // `seconds` == window size of one of the default options, `startSeconds` not specified.
      assert.deepEqual(
        findClosestTimeScale(
          defaultTimeScaleOptions,
          moment.duration(10, "minutes").asSeconds(),
        ),
        {
          ...defaultTimeScaleOptions["Past 10 Minutes"],
          key: "Past 10 Minutes",
        },
      );
      // `seconds` == window size of one of the default options, `startSeconds` not specified.
      assert.deepEqual(
        findClosestTimeScale(
          defaultTimeScaleOptions,
          moment.duration(14, "days").asSeconds(),
        ),
        {
          ...defaultTimeScaleOptions["Past 2 Weeks"],
          key: "Past 2 Weeks",
        },
      );
      // `seconds` == window size of one of the default options, `startSeconds` is now.
      assert.deepEqual(
        findClosestTimeScale(
          defaultTimeScaleOptions,
          defaultTimeScaleOptions["Past 10 Minutes"].windowSize.asSeconds(),
          moment().unix(),
        ),
        {
          ...defaultTimeScaleOptions["Past 10 Minutes"],
          key: "Past 10 Minutes",
        },
      );
      // `seconds` == window size of one of the default options, `startSeconds` is in the past.
      assert.deepEqual(
        findClosestTimeScale(
          defaultTimeScaleOptions,
          defaultTimeScaleOptions["Past 10 Minutes"].windowSize.asSeconds(),
          moment()
            .subtract(1, "day")
            .unix(),
        ),
        { ...defaultTimeScaleOptions["Past 10 Minutes"], key: "Custom" },
      );
    });
  });
});
