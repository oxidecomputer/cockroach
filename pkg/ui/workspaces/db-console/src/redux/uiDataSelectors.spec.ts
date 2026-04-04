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

import { assert } from "chai";
import { dismissReleaseNotesSignupForm } from "./uiDataSelectors";
import { UIData, UIDataStatus } from "src/redux/uiData";

describe("uiDataSelectors", () => {
  describe("dismissReleaseNotesSignupForm selector", () => {
    const selector = dismissReleaseNotesSignupForm.resultFunc;

    it("returns `false` if uiData status is VALID and has no data", () => {
      const uiData: UIData = {
        status: UIDataStatus.VALID,
        error: null,
        data: undefined,
      };
      assert.isFalse(selector(uiData));
    });

    it("returns `true` if uiData status is VALID and data = true", () => {
      const uiData: UIData = {
        status: UIDataStatus.VALID,
        error: null,
        data: true,
      };
      assert.isTrue(selector(uiData));
    });

    it("returns `true` if uiData status is UNINITIALIZED", () => {
      const uiData: UIData = {
        status: UIDataStatus.UNINITIALIZED,
        error: null,
        data: undefined,
      };
      assert.isTrue(selector(uiData));
    });

    it("returns `true` if uiData state is undefined", () => {
      assert.isTrue(selector(undefined));
    });
  });
});
