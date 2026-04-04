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

import React, { FunctionComponent } from "react";
import classnames from "classnames/bind";
import { Badge, BadgeIntent, FuzzyTime } from "@cockroachlabs/ui-components";

import {
  NotificationTypeProp,
  NotificationProps,
  NotificationSeverity,
} from "../Notifications";

import styles from "./notificationMessage.module.scss";

export type NotificationMessageProps = NotificationTypeProp & NotificationProps;

const cx = classnames.bind(styles);

const truncate = (string: string, length: number): string => {
  if (string.length <= length) return string;

  return `${string.slice(0, length)}…`;
};

const severityIntent = (s: NotificationSeverity): BadgeIntent => {
  const intentMap = {
    low: "neutral",
    info: "neutral",
    moderate: "info",
    critical: "info",
  };
  return intentMap[s] as BadgeIntent;
};

export const NotificationMessage: FunctionComponent<NotificationMessageProps> = ({
  id,
  description,
  read,
  severity,
  timestamp,
  title,
}) => {
  const time = new Date(timestamp);
  return (
    <section key={id} className={cx("notification-message", { unread: !read })}>
      <header className={cx("title")}>{title}</header>
      <Badge className={cx("severity")} intent={severityIntent(severity)}>
        {severity}
      </Badge>
      {description && (
        <div className={cx("description")}>{truncate(description, 120)}</div>
      )}
      <div className={cx("timestamp")}>
        <FuzzyTime timestamp={time} />
      </div>
    </section>
  );
};

export default NotificationMessage;
