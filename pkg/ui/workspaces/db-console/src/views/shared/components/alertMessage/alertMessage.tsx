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

import React from "react";
import { Alert, Icon } from "antd";
import { Link } from "react-router-dom";

import { AlertInfo, AlertLevel } from "src/redux/alerts";
import "./alertMessage.styl";

interface AlertMessageProps extends AlertInfo {
  autoClose: boolean;
  autoCloseTimeout: number;
  closable: boolean;
  dismiss(): void;
}

type AlertType = "success" | "info" | "warning" | "error";

const mapAlertLevelToType = (alertLevel: AlertLevel): AlertType => {
  switch (alertLevel) {
    case AlertLevel.SUCCESS:
      return "success";
    case AlertLevel.NOTIFICATION:
      return "info";
    case AlertLevel.WARNING:
      return "warning";
    case AlertLevel.CRITICAL:
      return "error";
    default:
      return "info";
  }
};

const getIconType = (alertLevel: AlertLevel): string => {
  switch (alertLevel) {
    case AlertLevel.SUCCESS:
      return "check-circle";
    case AlertLevel.NOTIFICATION:
      return "info-circle";
    case AlertLevel.WARNING:
      return "warning";
    case AlertLevel.CRITICAL:
      return "close-circle";
    default:
      return "info-circle";
  }
};

export class AlertMessage extends React.Component<AlertMessageProps> {
  static defaultProps = {
    closable: true,
    autoCloseTimeout: 6000,
  };

  timeoutHandler: number;

  componentDidMount() {
    const { autoClose, dismiss, autoCloseTimeout } = this.props;
    if (autoClose) {
      this.timeoutHandler = window.setTimeout(dismiss, autoCloseTimeout);
    }
  }

  componentWillUnmount() {
    clearTimeout(this.timeoutHandler);
  }

  render() {
    const { level, dismiss, link, title, text, closable } = this.props;

    let description: React.ReactNode = text;

    if (link) {
      description = (
        <Link to={link} onClick={dismiss}>
          {text}
        </Link>
      );
    }

    const type = mapAlertLevelToType(level);
    const iconType = getIconType(level);
    return (
      <Alert
        className="alert-massage"
        message={title}
        description={description}
        showIcon
        icon={
          <Icon
            type={iconType}
            theme="filled"
            className="alert-massage__icon"
          />
        }
        closable={closable}
        onClose={dismiss}
        closeText={
          closable && <div className="alert-massage__close-text">&times;</div>
        }
        type={type}
      />
    );
  }
}
