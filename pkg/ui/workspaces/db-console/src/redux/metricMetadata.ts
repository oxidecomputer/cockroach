// Copyright 2019 The Cockroach Authors.
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
import { MetricMetadataResponseMessage } from "src/util/api";

export type MetricsMetadata = MetricMetadataResponseMessage["metadata"];

// State selectors
const metricsMetadataStateSelector = (state: AdminUIState) =>
  state.cachedData.metricMetadata.data;

export const metricsMetadataSelector = createSelector(
  metricsMetadataStateSelector,
  (metricsMetadata): MetricsMetadata =>
    metricsMetadata ? metricsMetadata.metadata : undefined,
);
