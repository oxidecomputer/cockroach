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

import { createAction } from "@reduxjs/toolkit";
import { all, put, takeEvery } from "redux-saga/effects";

import { actions as terminateQueryActions } from "src/store/terminateQuery/terminateQuery.reducer";

export const notificationAction = createAction(
  "adminUI/notification",
  (type: NotificationType, text: string) => ({
    payload: {
      type,
      text,
    },
  }),
);

export enum NotificationType {
  Success = "success",
  Error = "error",
}

export type SendNotification = (
  type: NotificationType,
  message: string,
) => void;

export function* notifificationsSaga() {
  // ***************************** //
  // Terminate Query notifications //
  // ***************************** //
  yield all([
    takeEvery(terminateQueryActions.terminateSessionCompleted, function*() {
      yield put(
        notificationAction(NotificationType.Success, "Session cancelled."),
      );
    }),

    takeEvery(terminateQueryActions.terminateSessionFailed, function*() {
      yield put(
        notificationAction(
          NotificationType.Error,
          "There was an error cancelling the session",
        ),
      );
    }),

    takeEvery(terminateQueryActions.terminateQueryCompleted, function*() {
      yield put(
        notificationAction(NotificationType.Success, "Statement cancelled."),
      );
    }),

    takeEvery(terminateQueryActions.terminateQueryFailed, function*() {
      yield put(
        notificationAction(
          NotificationType.Error,
          "There was an error cancelling the statement.",
        ),
      );
    }),
  ]);
}
