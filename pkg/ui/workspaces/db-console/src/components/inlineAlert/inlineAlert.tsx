// Copyright 2020 The Cockroach Authors.
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

import React, { useMemo } from "react";
import classNames from "classnames/bind";

import styles from "./inlineAlert.module.styl";
import ErrorIcon from "assets/error-circle.svg";
import InfoIcon from "assets/info-filled-circle.svg";
import WarningIcon from "assets/warning.svg";

export type InlineAlertIntent = "info" | "error" | "warning";

const cn = classNames.bind(styles);

export interface InlineAlertProps {
  title: React.ReactNode;
  message?: React.ReactNode;
  intent?: InlineAlertIntent;
  className?: string;
  fullWidth?: boolean;
}

export const InlineAlert: React.FC<InlineAlertProps> = ({
  title,
  message,
  intent = "info",
  className,
  fullWidth,
}) => {
  const Icon = useMemo(() => {
    switch (intent) {
      case "error":
        return ErrorIcon;
      case "warning":
        return WarningIcon;
      case "info":
      default:
        return InfoIcon;
    }
  }, [intent]);

  return (
    <div
      className={cn(
        "root",
        `intent-${intent}`,
        { "full-width": fullWidth },
        className,
      )}
    >
      <div className={cn("icon-container")}>
        <img src={Icon} className={cn("icon")} />
      </div>
      <div className={cn("main-container")}>
        <div className={cn("title")}>{title}</div>
        <div className={cn("message")}>{message}</div>
      </div>
    </div>
  );
};
