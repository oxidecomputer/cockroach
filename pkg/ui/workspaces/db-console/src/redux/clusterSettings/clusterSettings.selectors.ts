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

import { createSelector } from "reselect";
import { AdminUIState } from "src/redux/state";
import { cockroach } from "src/js/protos";
import moment from "moment";
import { util } from "@cockroachlabs/cluster-ui";

export const selectClusterSettings = createSelector(
  (state: AdminUIState) => state.cachedData.settings?.data,
  (settings: cockroach.server.serverpb.SettingsResponse) =>
    settings?.key_values,
);

export const selectResolution10sStorageTTL = createSelector(
  selectClusterSettings,
  (settings): moment.Duration | undefined => {
    if (!settings) {
      return undefined;
    }
    const value = settings["timeseries.storage.resolution_10s.ttl"]?.value;
    return util.durationFromISO8601String(value);
  },
);

export const selectResolution30mStorageTTL = createSelector(
  selectClusterSettings,
  settings => {
    if (!settings) {
      return undefined;
    }
    const value = settings["timeseries.storage.resolution_30m.ttl"]?.value;
    return util.durationFromISO8601String(value);
  },
);

export const selectAutomaticStatsCollectionEnabled = createSelector(
  selectClusterSettings,
  (settings): boolean | undefined => {
    if (!settings) {
      return undefined;
    }
    const value = settings["sql.stats.automatic_collection.enabled"]?.value;
    return value === "true";
  },
);
