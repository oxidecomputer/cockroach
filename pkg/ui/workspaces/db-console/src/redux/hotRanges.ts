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

import { AdminUIState } from "src/redux/state";
import { createSelector } from "reselect";
import { cockroach } from "src/js/protos";

const hotRangesState = (state: AdminUIState) => state.cachedData.hotRanges;

export const hotRangesSelector = createSelector(hotRangesState, hotRanges =>
  Object.values(hotRanges?.data || {})
    .reduce<cockroach.server.serverpb.HotRangesResponseV2["ranges"]>(
      (acc, v) => [...acc, ...v.ranges],
      [],
    )
    // filter out ranges with 0 QPS
    .filter(v => v?.qps && v.qps > 0),
);

export const lastErrorSelector = createSelector(
  hotRangesState,
  hotRanges => hotRanges?.lastError,
);

export const isValidSelector = createSelector(
  hotRangesState,
  hotRanges => hotRanges?.valid,
);

export const lastSetAtSelector = createSelector(
  hotRangesState,
  hotRanges => hotRanges?.setAt,
);

export const isLoadingSelector = createSelector(
  hotRangesState,
  hotRanges => hotRanges?.inFlight,
);
