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

import classNames from "classnames/bind";
import statementsPageStyles from "src/statementsPage/statementsPage.module.scss";
import statementsTableStyles from "src/statementsTable/statementsTableContent.module.scss";
import sortedTableStyles from "src/sortedtable/sortedtable.module.scss";

const sortedTableCx = classNames.bind(sortedTableStyles);
const statementsTableCx = classNames.bind(statementsTableStyles);
const pageCx = classNames.bind(statementsPageStyles);

export const tableClasses = {
  containerClass: sortedTableCx("cl-table-container"),
  latencyClasses: {
    column: statementsTableCx("statements-table__col-latency"),
    barChart: {
      classes: {
        root: statementsTableCx("statements-table__col-latency--bar-chart"),
      },
    },
  },
};

export const statisticsClasses = {
  statistic: pageCx("cl-table-statistic"),
  countTitle: pageCx("cl-count-title"),
  lastCleared: pageCx("last-cleared-title"),
};
