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
import { Link } from "react-router-dom";
import { Tooltip } from "src/components";
import Job = cockroach.server.serverpb.IJobResponse;
import { cockroach } from "src/js/protos";

export class JobDescriptionCell extends React.PureComponent<{ job: Job }> {
  render() {
    // If a [SQL] job.statement exists, it means that job.description
    // is a human-readable message. Otherwise job.description is a SQL
    // statement.
    const job = this.props.job;
    const additionalStyle = job.statement ? "" : " jobs-table__cell--sql";
    const description =
      job.description && job.description.length > 425
        ? `${job.description.slice(0, 425)}...`
        : job.description;

    const cellContent = (
      <div className="jobs-table__cell--description">
        {job.statement || job.description || job.type}
      </div>
    );
    return (
      <Link className={`${additionalStyle}`} to={`jobs/${String(job.id)}`}>
        <div className="cl-table-link__tooltip">
          {description ? (
            <Tooltip
              arrowPointAtCenter
              placement="bottom"
              title={
                <pre
                  style={{ whiteSpace: "pre-wrap" }}
                  className="cl-table-link__description"
                >
                  {description}
                </pre>
              }
              overlayClassName="cl-table-link__statement-tooltip--fixed-width"
            >
              {cellContent}
            </Tooltip>
          ) : (
            cellContent
          )}
        </div>
      </Link>
    );
  }
}
