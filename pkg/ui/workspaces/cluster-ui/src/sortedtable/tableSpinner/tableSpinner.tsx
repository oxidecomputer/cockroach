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
import { Spinner } from "@cockroachlabs/icons";
import { Spin, Icon } from "antd";
import classNames from "classnames/bind";
import styles from "./tableSpinner.module.scss";

const cx = classNames.bind(styles);

interface TableSpinnerProps {
  loadingLabel: string;
}

export const TableSpinner = ({ loadingLabel }: TableSpinnerProps) => {
  const tableSpinnerClass = cx("table__loading");
  const spinClass = cx("table__loading--spin");
  const loadingLabelClass = cx("table__loading--label");

  return (
    <div className={tableSpinnerClass}>
      <Spin
        className={spinClass}
        indicator={<Icon component={Spinner} spin />}
      />
      {loadingLabel && (
        <span className={loadingLabelClass}>{loadingLabel}</span>
      )}
    </div>
  );
};
