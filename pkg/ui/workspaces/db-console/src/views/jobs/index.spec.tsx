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

import { assert } from "chai";
import moment from "moment";
import { cockroach } from "src/js/protos";
import { formatDuration } from ".";
import { JobsTable, JobsTableProps } from "src/views/jobs/index";
import {
  allJobsFixture,
  retryRunningJobFixture,
} from "src/views/jobs/jobsTable.fixture";
import { refreshJobs } from "src/redux/apiReducers";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import React from "react";
import { MemoryRouter } from "react-router-dom";
import * as H from "history";

import { expectPopperTooltipActivated } from "src/test-utils/tooltip";

import Job = cockroach.server.serverpb.IJobResponse;

const getMockJobsTableProps = (jobs: Array<Job>): JobsTableProps => {
  const history = H.createHashHistory();
  return {
    sort: { columnTitle: null, ascending: true },
    status: "",
    show: "50",
    type: 0,
    setSort: () => {},
    setStatus: () => {},
    setShow: () => {},
    setType: () => {},
    jobs: {
      data: {
        jobs: jobs,
      },
      inFlight: false,
      valid: true,
    },
    refreshJobs,
    location: history.location,
    history,
    match: {
      url: "",
      path: history.location.pathname,
      isExact: false,
      params: {},
    },
  };
};

describe("Jobs", () => {
  it("format duration", () => {
    assert.equal(formatDuration(moment.duration(0)), "00:00:00");
    assert.equal(formatDuration(moment.duration(5, "minutes")), "00:05:00");
    assert.equal(formatDuration(moment.duration(5, "hours")), "05:00:00");
    assert.equal(formatDuration(moment.duration(110, "hours")), "110:00:00");
    assert.equal(
      formatDuration(moment.duration(12345, "hours")),
      "12345:00:00",
    );
  });

  it("renders expected jobs table columns", () => {
    const { getByText } = render(
      <MemoryRouter>
        <JobsTable {...getMockJobsTableProps(allJobsFixture)} />
      </MemoryRouter>,
    );
    const expectedColumnTitles = [
      "Description",
      "Status",
      "Job ID",
      "User Name",
      "Creation Time (UTC)",
      "Last Execution Time (UTC)",
      "Execution Count",
    ];

    for (const columnTitle of expectedColumnTitles) {
      getByText(columnTitle);
    }
  });

  it("shows next execution time on hovering a retry status", async () => {
    const { getByText } = render(
      <MemoryRouter>
        <JobsTable {...getMockJobsTableProps([retryRunningJobFixture])} />
      </MemoryRouter>,
    );

    await waitFor(expectPopperTooltipActivated);
    userEvent.hover(getByText("retrying"));

    await waitFor(() =>
      screen.getByText("Next Execution Time", { exact: false }),
    );
  });
});
