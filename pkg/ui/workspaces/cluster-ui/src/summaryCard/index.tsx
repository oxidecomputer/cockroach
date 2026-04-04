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
import classnames from "classnames/bind";
import styles from "./summaryCard.module.scss";
import booleanSettingStyles from "../settings/booleanSetting.module.scss";
import { CircleFilled } from "src/icon";
import { Tooltip } from "antd";

interface ISummaryCardProps {
  children: React.ReactNode;
  className?: string;
  id?: string;
}

const cx = classnames.bind(styles);
const booleanSettingCx = classnames.bind(booleanSettingStyles);

// tslint:disable-next-line: variable-name
export const SummaryCard: React.FC<ISummaryCardProps> = ({
  children,
  className = "",
  id,
}) => (
  <div className={`${cx("summary--card")} ${className}`} id={id}>
    {children}
  </div>
);

interface ISummaryCardItemProps {
  label: React.ReactNode;
  value: React.ReactNode;
  className?: string;
}

interface ISummaryCardItemBoolSettingProps extends ISummaryCardItemProps {
  toolTipText: JSX.Element;
}

export const SummaryCardItem: React.FC<ISummaryCardItemProps> = ({
  label,
  value,
  className = "",
}) => (
  <div className={cx("summary--card__item", className)}>
    <h4 className={cx("summary--card__item--label")}>{label}</h4>
    <p className={cx("summary--card__item--value")}>{value}</p>
  </div>
);

export const SummaryCardItemBoolSetting: React.FC<ISummaryCardItemBoolSettingProps> = ({
  label,
  value,
  toolTipText,
  className,
}) => {
  const boolValue = value ? "Enabled" : "Disabled";
  const boolClass = value
    ? "bool-setting-icon__enabled"
    : "bool-setting-icon__disabled";

  return (
    <div className={cx("summary--card__item", className)}>
      <h4 className={cx("summary--card__item--label")}>{label}</h4>
      <p className={cx("summary--card__item--value")}>
        <CircleFilled className={booleanSettingCx(boolClass)} />
        <Tooltip
          placement="bottom"
          title={toolTipText}
          className={cx("crl-hover-text__dashed-underline")}
        >
          {boolValue}
        </Tooltip>
      </p>
    </div>
  );
};
