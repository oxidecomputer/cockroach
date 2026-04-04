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

import { cockroach } from "src/js/protos";
import Job = cockroach.server.serverpb.IJobResponse;
import { BadgeStatus } from "src/components";

export enum JobStatusVisual {
  BadgeOnly,
  BadgeWithDuration,
  ProgressBarWithDuration,
  BadgeWithMessage,
  BadgeWithErrorMessage,
  BadgeWithRetrying,
}

export function jobToVisual(job: Job): JobStatusVisual {
  if (job.type === "CHANGEFEED") {
    return JobStatusVisual.BadgeOnly;
  }
  switch (job.status) {
    case JOB_STATUS_SUCCEEDED:
      return JobStatusVisual.BadgeWithDuration;
    case JOB_STATUS_FAILED:
      return JobStatusVisual.BadgeWithErrorMessage;
    case JOB_STATUS_RUNNING:
      return JobStatusVisual.ProgressBarWithDuration;
    case JOB_STATUS_RETRY_RUNNING:
      return JobStatusVisual.ProgressBarWithDuration;
    case JOB_STATUS_PENDING:
      return JobStatusVisual.BadgeWithMessage;
    case JOB_STATUS_RETRY_REVERTING:
      return JobStatusVisual.BadgeWithRetrying;
    case JOB_STATUS_CANCELED:
    case JOB_STATUS_CANCEL_REQUESTED:
    case JOB_STATUS_PAUSED:
      return job.error == ""
        ? JobStatusVisual.BadgeOnly
        : JobStatusVisual.BadgeWithErrorMessage;
    case JOB_STATUS_PAUSE_REQUESTED:
    case JOB_STATUS_REVERTING:
    default:
      return JobStatusVisual.BadgeOnly;
  }
}

export const JOB_STATUS_SUCCEEDED = "succeeded";
export const JOB_STATUS_FAILED = "failed";
export const JOB_STATUS_CANCELED = "canceled";
export const JOB_STATUS_CANCEL_REQUESTED = "cancel-requested";
export const JOB_STATUS_PAUSED = "paused";
export const JOB_STATUS_PAUSE_REQUESTED = "paused-requested";
export const JOB_STATUS_RUNNING = "running";
export const JOB_STATUS_RETRY_RUNNING = "retry-running";
export const JOB_STATUS_PENDING = "pending";
export const JOB_STATUS_REVERTING = "reverting";
export const JOB_STATUS_REVERT_FAILED = "revert-failed";
export const JOB_STATUS_RETRY_REVERTING = "retry-reverting";

export function isRetrying(status: string): boolean {
  return [JOB_STATUS_RETRY_RUNNING, JOB_STATUS_RETRY_REVERTING].includes(
    status,
  );
}
export function isRunning(status: string): boolean {
  return [JOB_STATUS_RUNNING, JOB_STATUS_RETRY_RUNNING].includes(status);
}

export const statusOptions = [
  { value: "", label: "All" },
  { value: "succeeded", label: "Succeeded" },
  { value: "failed", label: "Failed" },
  { value: "paused", label: "Paused" },
  { value: "canceled", label: "Canceled" },
  { value: "running", label: "Running" },
  { value: "pending", label: "Pending" },
  { value: "reverting", label: "Reverting" },
  { value: "retrying", label: "Retrying" },
];

export function jobHasOneOfStatuses(job: Job, ...statuses: string[]) {
  return statuses.indexOf(job.status) !== -1;
}

export const jobStatusToBadgeStatus = (status: string): BadgeStatus => {
  switch (status) {
    case JOB_STATUS_SUCCEEDED:
      return "success";
    case JOB_STATUS_FAILED:
      return "danger";
    case JOB_STATUS_REVERT_FAILED:
      return "danger";
    case JOB_STATUS_RUNNING:
      return "info";
    case JOB_STATUS_PENDING:
      return "warning";
    case JOB_STATUS_CANCELED:
    case JOB_STATUS_CANCEL_REQUESTED:
    case JOB_STATUS_PAUSED:
    case JOB_STATUS_PAUSE_REQUESTED:
    case JOB_STATUS_REVERTING:
    case JOB_STATUS_RETRY_REVERTING:
    default:
      return "default";
  }
};
export const jobStatusToBadgeText = (status: string): string => {
  switch (status) {
    case JOB_STATUS_RETRY_REVERTING:
      return JOB_STATUS_REVERTING;
    case JOB_STATUS_RETRY_RUNNING:
      return JOB_STATUS_RUNNING;
    default:
      return status;
  }
};
