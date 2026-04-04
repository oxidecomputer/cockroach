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
import { cockroach } from "src/js/protos";
import { JobStatus } from "./jobStatus";
import { isRetrying } from "src/views/jobs/jobStatusOptions";
import { util } from "@cockroachlabs/cluster-ui";
import { DATE_FORMAT_24_UTC } from "src/util/format";
import { Tooltip } from "@cockroachlabs/ui-components";
import Job = cockroach.server.serverpb.IJobResponse;

export interface JobStatusCellProps {
  job: Job;
  lineWidth?: number;
  compact?: boolean;
}

export const JobStatusCell: React.FC<JobStatusCellProps> = ({
  job,
  lineWidth,
  compact = false,
}) => {
  const jobStatus = (
    <JobStatus job={job} lineWidth={lineWidth} compact={compact} />
  );
  if (isRetrying(job.status)) {
    return (
      <Tooltip
        placement="bottom"
        style="tableTitle"
        content={
          <>
            Next Execution Time:
            <br />
            {util.TimestampToMoment(job.next_run).format(DATE_FORMAT_24_UTC)}
          </>
        }
      >
        {jobStatus}
      </Tooltip>
    );
  }
  return jobStatus;
};
