// Copyright 2021 The Cockroach Authors.
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

import { Filters, getFiltersFromQueryString } from "./filter";

describe("Test filter functions", (): void => {
  describe("Test get filters from query string", (): void => {
    it("no values on query string", (): void => {
      const expectedFilters: Filters = {
        app: "",
        timeNumber: "0",
        timeUnit: "seconds",
        fullScan: false,
        sqlType: "",
        database: "",
        regions: "",
        nodes: "",
      };
      const resultFilters = getFiltersFromQueryString("");
      expect(resultFilters).toEqual(expectedFilters);
    });
  });

  it("different values from default values on query string", (): void => {
    const expectedFilters: Filters = {
      app: "$ internal",
      timeNumber: "1",
      timeUnit: "milliseconds",
      fullScan: true,
      sqlType: "DML",
      database: "movr",
      regions: "us-central",
      nodes: "n1,n2",
    };
    const resultFilters = getFiltersFromQueryString(
      "app=%24+internal&timeNumber=1&timeUnit=milliseconds&fullScan=true&sqlType=DML&database=movr&regions=us-central&nodes=n1,n2",
    );
    expect(resultFilters).toEqual(expectedFilters);
  });

  it("testing boolean with full scan = true", (): void => {
    const expectedFilters: Filters = {
      app: "",
      timeNumber: "0",
      timeUnit: "seconds",
      fullScan: true,
      sqlType: "",
      database: "",
      regions: "",
      nodes: "",
    };
    const resultFilters = getFiltersFromQueryString("fullScan=true");
    expect(resultFilters).toEqual(expectedFilters);
  });

  it("testing boolean with full scan = false", (): void => {
    const expectedFilters: Filters = {
      app: "",
      timeNumber: "0",
      timeUnit: "seconds",
      fullScan: false,
      sqlType: "",
      database: "",
      regions: "",
      nodes: "",
    };
    const resultFilters = getFiltersFromQueryString("fullScan=false");
    expect(resultFilters).toEqual(expectedFilters);
  });
});
