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
import {
  jobStatusToBadgeStatus,
  jobStatusToBadgeText,
} from "src/views/jobs/jobStatusOptions";
import Job = cockroach.server.serverpb.IJobResponse;
import { cockroach } from "src/js/protos";
import { Badge } from "src/components";
import { Line } from "rc-progress";
import { ColorIntentInfo3 } from "@cockroachlabs/design-tokens";

export class JobStatusBadge extends React.PureComponent<{ jobStatus: string }> {
  render() {
    const jobStatus = this.props.jobStatus;
    const badgeStatus = jobStatusToBadgeStatus(jobStatus);
    const badgeText = jobStatusToBadgeText(jobStatus);
    return <Badge status={badgeStatus} text={badgeText} />;
  }
}

export class RetryingStatusBadge extends React.PureComponent {
  render() {
    return <Badge status="warning" text="retrying" />;
  }
}

export class ProgressBar extends React.PureComponent<{
  job: Job;
  lineWidth: number;
  showPercentage: boolean;
}> {
  render() {
    const percent = this.props.job.fraction_completed * 100;
    return (
      <div className="jobs-table__progress">
        <Line
          percent={percent}
          strokeWidth={this.props.lineWidth}
          trailWidth={this.props.lineWidth}
          strokeColor={ColorIntentInfo3}
          trailColor="#d6dbe7"
          className="jobs-table__progress-bar"
        />
        {this.props.showPercentage ? (
          <div
            className="jobs-table__status--percentage"
            title={percent.toFixed(3) + "%"}
          >
            {percent.toFixed(1) + "%"}
          </div>
        ) : null}
      </div>
    );
  }
}
