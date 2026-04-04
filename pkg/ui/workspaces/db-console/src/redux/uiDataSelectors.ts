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

import { createSelector } from "reselect";

import { AdminUIState } from "src/redux/state";
import {
  RELEASE_NOTES_SIGNUP_DISMISSED_KEY,
  UIDataStatus,
} from "src/redux/uiData";

export const dismissReleaseNotesSignupForm = createSelector(
  (state: AdminUIState) => state.uiData[RELEASE_NOTES_SIGNUP_DISMISSED_KEY],
  hideFormData => {
    // Do not show subscription form if data is not initialized yet.
    // It avoids form flickering in case value is set to `false` (hide form) and it
    // is shown for a moment before response is received back.
    if (!hideFormData) {
      return true;
    }
    if (hideFormData.status === UIDataStatus.VALID) {
      // If data is successfully loaded and have no values,
      // return default `false` value (do not hide subscription form)
      if (hideFormData?.data === undefined) {
        return false;
      }
      return hideFormData?.data;
    }
    // Do not show subscription form if request is loading
    return true;
  },
);
