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

import React from "react";
import { statisticsClasses } from "../transactionsPage/transactionsPageClasses";
import { ISortedTablePagination } from "../sortedtable";
import { Button } from "src/button";
import { ResultsPerPageLabel } from "src/pagination";
import classNames from "classnames/bind";
import timeScaleStyles from "../timeScaleDropdown/timeScale.module.scss";

const { statistic, countTitle } = statisticsClasses;
const timeScaleStylesCx = classNames.bind(timeScaleStyles);

interface TableStatistics {
  pagination: ISortedTablePagination;
  totalCount: number;
  arrayItemName: string;
  activeFilters: number;
  search?: string;
  onClearFilters?: () => void;
  period?: string;
}

// TableStatistics shows statistics about the results being
// displayed on the table, including total results count,
// how many are currently being displayed on the page and
// if there is an active filter.
// This component has also a clear filter option.
export const TableStatistics: React.FC<TableStatistics> = ({
  pagination,
  totalCount,
  search,
  arrayItemName,
  onClearFilters,
  activeFilters,
  period,
}) => {
  const periodLabel = (
    <>
      &nbsp;&nbsp;&nbsp;| &nbsp;
      <p className={timeScaleStylesCx("time-label")}>
        Showing aggregated stats from{" "}
        <span className={timeScaleStylesCx("bold")}>{period}</span>
      </p>
    </>
  );

  const resultsPerPageLabel = (
    <>
      <ResultsPerPageLabel
        pagination={{ ...pagination, total: totalCount }}
        pageName={arrayItemName}
        search={search}
      />
      {period && periodLabel}
    </>
  );

  const resultsCountAndClear = (
    <>
      {totalCount} {totalCount === 1 ? "result" : "results"}
      {period && periodLabel}
      &nbsp;&nbsp;&nbsp;| &nbsp;
      <Button onClick={() => onClearFilters()} type="flat" size="small">
        clear filter
      </Button>
    </>
  );

  return (
    <div className={statistic}>
      <h4 className={countTitle}>
        {activeFilters ? resultsCountAndClear : resultsPerPageLabel}
      </h4>
    </div>
  );
};
