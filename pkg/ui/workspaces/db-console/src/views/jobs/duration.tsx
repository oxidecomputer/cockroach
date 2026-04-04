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
import { util } from "@cockroachlabs/cluster-ui";
import {
  isRunning,
  JOB_STATUS_SUCCEEDED,
} from "src/views/jobs/jobStatusOptions";
import { formatDuration } from "src/views/jobs/index";
import moment from "moment";
import Job = cockroach.server.serverpb.IJobResponse;
import { cockroach } from "src/js/protos";

export class Duration extends React.PureComponent<{
  job: Job;
  className?: string;
}> {
  render() {
    const { job, className } = this.props;
    // Parse timestamp to default value NULL instead of Date.now.
    // Conversion dates to Date.now causes trailing dates and constant
    // duration increase even when job is finished.
    const startedAt = util.TimestampToMoment(job.started, null);
    const modifiedAt = util.TimestampToMoment(job.modified, null);
    const finishedAt = util.TimestampToMoment(job.finished, null);

    if (isRunning(job.status)) {
      const fractionCompleted = job.fraction_completed;
      if (fractionCompleted > 0) {
        const duration = modifiedAt.diff(startedAt);
        const remaining = duration / fractionCompleted - duration;
        return (
          <span className={className}>
            {formatDuration(moment.duration(remaining)) + " remaining"}
          </span>
        );
      }
      return null;
    } else if (job.status == JOB_STATUS_SUCCEEDED) {
      return (
        <span className={className}>
          {"Duration: " +
            formatDuration(moment.duration(finishedAt.diff(startedAt)))}
        </span>
      );
    }
    return null;
  }
}
