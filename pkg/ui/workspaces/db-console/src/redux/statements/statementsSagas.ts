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

import { all, call, put, takeEvery, takeLatest } from "redux-saga/effects";
import { PayloadAction } from "src/interfaces/action";

import {
  cancelStatementDiagnosticsReport,
  CancelStatementDiagnosticsReportResponseMessage,
  createStatementDiagnosticsReport,
} from "src/util/api";
import {
  CREATE_STATEMENT_DIAGNOSTICS_REPORT,
  CreateStatementDiagnosticsReportPayload,
  createStatementDiagnosticsReportCompleteAction,
  createStatementDiagnosticsReportFailedAction,
  SET_GLOBAL_TIME_SCALE,
  CANCEL_STATEMENT_DIAGNOSTICS_REPORT,
  cancelStatementDiagnosticsReportCompleteAction,
  cancelStatementDiagnosticsReportFailedAction,
  CancelStatementDiagnosticsReportPayload,
} from "./statementsActions";
import { cockroach } from "src/js/protos";
import CreateStatementDiagnosticsReportRequest = cockroach.server.serverpb.CreateStatementDiagnosticsReportRequest;
import CancelStatementDiagnosticsReportRequest = cockroach.server.serverpb.CancelStatementDiagnosticsReportRequest;
import {
  invalidateStatementDiagnosticsRequests,
  refreshStatementDiagnosticsRequests,
} from "src/redux/apiReducers";
import {
  createStatementDiagnosticsAlertLocalSetting,
  cancelStatementDiagnosticsAlertLocalSetting,
} from "src/redux/alerts";
import { TimeScale } from "@cockroachlabs/cluster-ui";
import { setTimeScale } from "src/redux/timeScale";

export function* createDiagnosticsReportSaga(
  action: PayloadAction<CreateStatementDiagnosticsReportPayload>,
) {
  const { statementFingerprint, minExecLatency, expiresAfter } = action.payload;
  const createDiagnosticsReportRequest = new CreateStatementDiagnosticsReportRequest(
    {
      statement_fingerprint: statementFingerprint,
      min_execution_latency: minExecLatency,
      expires_after: expiresAfter,
    },
  );
  try {
    yield call(
      createStatementDiagnosticsReport,
      createDiagnosticsReportRequest,
    );
    yield put(createStatementDiagnosticsReportCompleteAction());
    yield put(invalidateStatementDiagnosticsRequests());
    // PUT expects action with `type` field which isn't defined in `refresh` ThunkAction interface
    yield put(refreshStatementDiagnosticsRequests() as any);
    // Stop showing the "cancel statement" alert if it is currently showing
    // (i.e. accidental cancel, then immediate activate).
    yield put(
      cancelStatementDiagnosticsAlertLocalSetting.set({
        show: false,
      }),
    );
    yield put(
      createStatementDiagnosticsAlertLocalSetting.set({
        show: true,
        status: "SUCCESS",
      }),
    );
  } catch (e) {
    yield put(createStatementDiagnosticsReportFailedAction());
    // Stop showing the "cancel statement" alert if it is currently showing
    // (i.e. accidental cancel, then immediate activate).
    yield put(
      cancelStatementDiagnosticsAlertLocalSetting.set({
        show: false,
      }),
    );
    yield put(
      createStatementDiagnosticsAlertLocalSetting.set({
        show: true,
        status: "FAILED",
      }),
    );
  }
}

export function* cancelDiagnosticsReportSaga(
  action: PayloadAction<CancelStatementDiagnosticsReportPayload>,
) {
  const { requestID } = action.payload;
  const cancelDiagnosticsReportRequest = new CancelStatementDiagnosticsReportRequest(
    {
      request_id: requestID,
    },
  );
  try {
    const response: CancelStatementDiagnosticsReportResponseMessage = yield call(
      cancelStatementDiagnosticsReport,
      cancelDiagnosticsReportRequest,
    );

    if (response.error !== "") {
      throw response.error;
    }

    yield put(cancelStatementDiagnosticsReportCompleteAction());

    yield put(invalidateStatementDiagnosticsRequests());
    // PUT expects action with `type` field which isn't defined in `refresh` ThunkAction interface.
    yield put(refreshStatementDiagnosticsRequests() as any);

    // Stop showing the "create statement" alert if it is currently showing
    // (i.e. accidental activate, then immediate cancel).
    yield put(
      createStatementDiagnosticsAlertLocalSetting.set({
        show: false,
      }),
    );
    yield put(
      cancelStatementDiagnosticsAlertLocalSetting.set({
        show: true,
        status: "SUCCESS",
      }),
    );
  } catch (e) {
    yield put(cancelStatementDiagnosticsReportFailedAction());
    // Stop showing the "create statement" alert if it is currently showing
    // (i.e. accidental activate, then immediate cancel).
    yield put(
      createStatementDiagnosticsAlertLocalSetting.set({
        show: false,
      }),
    );
    yield put(
      cancelStatementDiagnosticsAlertLocalSetting.set({
        show: true,
        status: "FAILED",
      }),
    );
  }
}

export function* setCombinedStatementsTimeScaleSaga(
  action: PayloadAction<TimeScale>,
) {
  const ts = action.payload;

  yield put(setTimeScale(ts));
}

export function* statementsSaga() {
  yield all([
    takeEvery(CREATE_STATEMENT_DIAGNOSTICS_REPORT, createDiagnosticsReportSaga),
    takeEvery(CANCEL_STATEMENT_DIAGNOSTICS_REPORT, cancelDiagnosticsReportSaga),
    takeLatest(SET_GLOBAL_TIME_SCALE, setCombinedStatementsTimeScaleSaga),
  ]);
}
