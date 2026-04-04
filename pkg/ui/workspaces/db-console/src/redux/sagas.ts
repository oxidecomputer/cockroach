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

import { all, fork } from "redux-saga/effects";

import { queryMetricsSaga } from "./metrics";
import { localSettingsSaga } from "./localsettings";
import { statementsSaga } from "./statements";
import { sessionsSaga } from "./sessions";
import { sqlStatsSaga } from "./sqlStats";
import { indexUsageStatsSaga } from "./indexUsageStats";
import { timeScaleSaga } from "src/redux/timeScale";

export default function* rootSaga() {
  yield all([
    fork(queryMetricsSaga),
    fork(localSettingsSaga),
    fork(statementsSaga),
    fork(sessionsSaga),
    fork(sqlStatsSaga),
    fork(indexUsageStatsSaga),
    fork(timeScaleSaga),
  ]);
}
