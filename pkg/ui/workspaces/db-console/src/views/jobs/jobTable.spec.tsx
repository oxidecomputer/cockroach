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
import { shallow } from "enzyme";
import { assert } from "chai";
import { JobTable, JobTableProps } from "src/views/jobs/jobTable";

import "src/enzymeInit";

describe("<JobTable>", () => {
  it("should reset page to 1 after job list prop changes", () => {
    const jobTableProps: JobTableProps = {
      sort: { columnTitle: null, ascending: true },
      setSort: () => {},
      jobs: {
        data: { jobs: [{}, {}, {}, {}] },
        inFlight: false,
        valid: true,
      },
      current: 2,
      pageSize: 2,
      isUsedFilter: true,
    };
    const jobTable = shallow<JobTable>(
      <JobTable
        jobs={jobTableProps.jobs}
        sort={jobTableProps.sort}
        setSort={jobTableProps.setSort}
        current={jobTableProps.current}
        pageSize={jobTableProps.pageSize}
        isUsedFilter={jobTableProps.isUsedFilter}
      />,
    );
    assert.equal(jobTable.state().pagination.current, 2);
    jobTable.setProps({
      jobs: { data: { jobs: [{}, {}] }, inFlight: false, valid: true },
    });
    assert.equal(jobTable.state().pagination.current, 1);
  });
});
