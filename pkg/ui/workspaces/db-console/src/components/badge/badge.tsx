// Copyright 2019 The Cockroach Authors.
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
import classNames from "classnames/bind";

import styles from "./badge.module.styl";

export type BadgeStatus = "success" | "danger" | "default" | "info" | "warning";

export interface BadgeProps {
  text: React.ReactNode;
  size?: "small" | "medium" | "large";
  status?: BadgeStatus;
  icon?: React.ReactNode;
  iconPosition?: "left" | "right";
}

Badge.defaultProps = {
  size: "medium",
  status: "default",
};

const cx = classNames.bind(styles);

export function Badge(props: BadgeProps) {
  const { size, status, icon, iconPosition, text } = props;
  const classes = cx("badge", `badge--size-${size}`, `badge--status-${status}`);
  const iconClasses = cx(
    "badge__icon",
    `badge__icon--position-${iconPosition || "left"}`,
  );
  return (
    <div className={classes}>
      {icon && <div className={iconClasses}>{icon}</div>}
      <div className={cx("badge__text", "badge__text--no-wrap")}>{text}</div>
    </div>
  );
}
