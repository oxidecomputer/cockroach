// Copyright 2023 The Cockroach Authors.
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
import styles from "./searchCriteria.module.scss";
import {
  Button,
  commonStyles,
  PageConfig,
  PageConfigItem,
  selectCustomStyles,
  TimeScale,
  timeScale1hMinOptions,
  TimeScaleDropdown,
} from "src";
import { applyBtn } from "../queryFilter/filterClasses";
import Select from "react-select";
import { limitOptions } from "../util/sqlActivityConstants";
import { SqlStatsSortType } from "src/api/statementsApi";
const cx = classNames.bind(styles);

type SortOption = {
  label: string;
  value: SqlStatsSortType;
};
export interface SearchCriteriaProps {
  sortOptions: SortOption[];
  currentScale: TimeScale;
  topValue: number;
  byValue: SqlStatsSortType;
  onChangeTimeScale: (ts: TimeScale) => void;
  onChangeTop: (top: number) => void;
  onChangeBy: (by: SqlStatsSortType) => void;
  onApply: () => void;
}

export function SearchCriteria(props: SearchCriteriaProps): React.ReactElement {
  const {
    topValue,
    byValue,
    currentScale,
    onChangeTop,
    onChangeBy,
    onChangeTimeScale,
    sortOptions,
  } = props;
  const customStyles = { ...selectCustomStyles };
  customStyles.indicatorSeparator = (provided: any) => ({
    ...provided,
    display: "none",
  });

  const customStylesTop = { ...customStyles };
  customStylesTop.container = (provided: any) => ({
    ...provided,
    width: "80px",
    border: "none",
  });

  const customStylesBy = { ...customStyles };
  customStylesBy.container = (provided: any) => ({
    ...provided,
    width: "170px",
    border: "none",
  });

  return (
    <div className={cx("search-area")}>
      <h5 className={commonStyles("base-heading")}>Search Criteria</h5>
      <PageConfig className={cx("top-area")}>
        <PageConfigItem>
          <label>
            <span className={cx("label")}>Top</span>
            <Select
              options={limitOptions}
              value={limitOptions.filter(top => top.value === topValue)}
              onChange={event => onChangeTop(event.value)}
              styles={customStylesTop}
            />
          </label>
        </PageConfigItem>
        <PageConfigItem>
          <label>
            <span className={cx("label")}>By</span>
            <Select
              options={sortOptions}
              value={sortOptions.filter(
                (top: SortOption) => top.value === byValue,
              )}
              onChange={event => onChangeBy(event.value as SqlStatsSortType)}
              styles={customStylesBy}
            />
          </label>
        </PageConfigItem>
        <PageConfigItem>
          <label>
            <span className={cx("label")}>Time Range</span>
            <TimeScaleDropdown
              options={timeScale1hMinOptions}
              currentScale={currentScale}
              setTimeScale={onChangeTimeScale}
              className={cx("timescale-small")}
            />
          </label>
        </PageConfigItem>
        <PageConfigItem>
          <Button
            className={`${applyBtn.btn} ${cx("margin-top-btn")}`}
            textAlign="center"
            onClick={props.onApply}
          >
            Apply
          </Button>
        </PageConfigItem>
      </PageConfig>
    </div>
  );
}
