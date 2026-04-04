// Copyright 2018 The Cockroach Authors.
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

import { assert } from "chai";
import { defaultTimeScaleOptions } from "@cockroachlabs/cluster-ui";
import * as timeScale from "./timeScale";
import moment from "moment";
import sinon from "sinon";

const sandbox = sinon.createSandbox();

describe("time scale reducer", function() {
  beforeEach(() => {
    sandbox.stub(sessionStorage, "getItem").returns(null);
  });

  afterEach(() => {
    sandbox.restore();
  });

  describe("actions", function() {
    it("should create the correct SET_METRICS_MOVING_WINDOW action to set the current time window", function() {
      const start = moment();
      const end = start.add(10, "s");
      const expectedSetting = {
        type: timeScale.SET_METRICS_MOVING_WINDOW,
        payload: {
          start,
          end,
        },
      };
      assert.deepEqual(
        timeScale.setMetricsMovingWindow({ start, end }),
        expectedSetting,
      );
    });

    it("should create the correct SET_SCALE action to set time window settings", function() {
      const payload: timeScale.TimeScale = {
        windowSize: moment.duration(10, "s"),
        windowValid: moment.duration(10, "s"),
        sampleSize: moment.duration(10, "s"),
        fixedWindowEnd: false,
      };
      assert.deepEqual(timeScale.setTimeScale(payload), {
        type: timeScale.SET_SCALE,
        payload,
      });
    });
  });

  describe("reducer", () => {
    it("should have the correct default value.", () => {
      assert.deepEqual(
        timeScale.timeScaleReducer(undefined, { type: "unknown" }),
        new timeScale.TimeScaleState(),
      );
      assert.deepEqual(new timeScale.TimeScaleState().scale, {
        ...defaultTimeScaleOptions["Past 10 Minutes"],
        key: "Past 10 Minutes",
        fixedWindowEnd: false,
      });
    });

    describe("setMetricsMovingWindow", () => {
      const start = moment();
      const end = start.add(10, "s");
      it("should correctly overwrite previous value", () => {
        const expected = new timeScale.TimeScaleState();
        expected.metricsTime.currentWindow = {
          start,
          end,
        };
        expected.metricsTime.shouldUpdateMetricsWindowFromScale = false;
        assert.deepEqual(
          timeScale.timeScaleReducer(
            undefined,
            timeScale.setMetricsMovingWindow({ start, end }),
          ),
          expected,
        );
      });
    });

    describe("setTimeScale", () => {
      const newSize = moment.duration(1, "h");
      const newValid = moment.duration(1, "m");
      const newSample = moment.duration(1, "m");
      it("should correctly overwrite previous value", () => {
        const expected = new timeScale.TimeScaleState();
        expected.scale = {
          windowSize: newSize,
          windowValid: newValid,
          sampleSize: newSample,
          fixedWindowEnd: false,
        };
        expected.metricsTime.shouldUpdateMetricsWindowFromScale = true;
        assert.deepEqual(
          timeScale.timeScaleReducer(
            undefined,
            timeScale.setTimeScale({
              windowSize: newSize,
              windowValid: newValid,
              sampleSize: newSample,
              fixedWindowEnd: false,
            }),
          ),
          expected,
        );
      });
    });
  });
});
