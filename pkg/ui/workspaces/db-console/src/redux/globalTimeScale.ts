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

import { LocalSetting } from "./localsettings";
import { AdminUIState } from "./state";
import { TimeScale, defaultTimeScaleSelected } from "@cockroachlabs/cluster-ui";

const localSettingsSelector = (state: AdminUIState) => state.localSettings;

export const globalTimeScaleLocalSetting = new LocalSetting<
  AdminUIState,
  TimeScale
>("timeScale/SQLActivity", localSettingsSelector, defaultTimeScaleSelected);
