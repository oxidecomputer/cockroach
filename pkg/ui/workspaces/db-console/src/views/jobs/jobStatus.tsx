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

import React from "react";
import classNames from "classnames/bind";
import {
  JobStatusBadge,
  RetryingStatusBadge,
  ProgressBar,
} from "src/views/jobs/progressBar";
import { Duration } from "src/views/jobs/duration";
import Job = cockroach.server.serverpb.IJobResponse;
import { cockroach } from "src/js/protos";
import {
  JobStatusVisual,
  jobToVisual,
  isRetrying,
} from "src/views/jobs/jobStatusOptions";
import { InlineAlert } from "src/components";
import styles from "./jobStatus.module.styl";

export interface JobStatusProps {
  job: Job;
  lineWidth?: number;
  compact?: boolean;
}

const cn = classNames.bind(styles);

export const JobStatus: React.FC<JobStatusProps> = ({
  job,
  compact,
  lineWidth,
}) => {
  const visualType = jobToVisual(job);

  switch (visualType) {
    case JobStatusVisual.BadgeOnly:
      return <JobStatusBadge jobStatus={job.status} />;
    case JobStatusVisual.BadgeWithDuration:
      return (
        <div>
          <JobStatusBadge jobStatus={job.status} />
          <Duration job={job} className="jobs-table__duration" />
        </div>
      );
    case JobStatusVisual.ProgressBarWithDuration: {
      const jobIsRetrying = isRetrying(job.status);
      return (
        <div>
          <ProgressBar
            job={job}
            lineWidth={lineWidth || 11}
            showPercentage={true}
          />
          <Duration job={job} className={cn("jobs-table__duration")} />
          {jobIsRetrying && <RetryingStatusBadge />}
          {job.running_status && (
            <div className="jobs-table__running-status">
              {job.running_status}
            </div>
          )}
        </div>
      );
    }
    case JobStatusVisual.BadgeWithMessage:
      return (
        <div>
          <JobStatusBadge jobStatus={job.status} />
          <span className="jobs-table__running-status">
            {job.running_status}
          </span>
        </div>
      );
    case JobStatusVisual.BadgeWithErrorMessage:
      return (
        <div>
          <JobStatusBadge jobStatus={job.status} />
          {!compact && (
            <InlineAlert
              title={job.error}
              intent="error"
              className={cn("inline-message")}
            />
          )}
        </div>
      );
    case JobStatusVisual.BadgeWithRetrying:
      return (
        <div className="jobs-table__two-statuses">
          <JobStatusBadge jobStatus={job.status} />
          <RetryingStatusBadge />
        </div>
      );
    default:
      return <JobStatusBadge jobStatus={job.status} />;
  }
};
