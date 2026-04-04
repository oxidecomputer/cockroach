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

import * as React from "react";
import { CircleFilled } from "src/icon";
import { Tooltip } from "antd";
import classNames from "classnames/bind";
import styles from "./booleanSetting.module.scss";

const cx = classNames.bind(styles);

export interface BooleanSettingProps {
  text: string;
  enabled: boolean;
  tooltipText: JSX.Element;
}

export function BooleanSetting(props: BooleanSettingProps) {
  const { text, enabled, tooltipText } = props;
  const label = enabled ? "enabled" : "disabled";
  const boolClass = enabled
    ? "bool-setting-icon__enabled"
    : "bool-setting-icon__disabled";
  return (
    <div>
      <CircleFilled className={cx(boolClass)} />
      <Tooltip
        placement="bottom"
        title={tooltipText}
        className={cx("crl-hover-text__dashed-underline")}
      >
        {text} - {label}
      </Tooltip>
    </div>
  );
}
