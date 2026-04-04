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
import classNames from "classnames/bind";
import styles from "./sqlActivity.module.scss";

const cx = classNames.bind(styles);

interface SQLActivityErrorProps {
  statsType: string;
  timeout?: boolean;
}

const LoadingError: React.FC<SQLActivityErrorProps> = props => {
  const error = props.timeout ? "a timeout" : "an unexpected error";
  return (
    <div className={cx("row")}>
      <span>{`This page had ${error} while loading ${props.statsType}.`}</span>
      &nbsp;
      <a
        className={cx("action")}
        onClick={() => {
          window.location.reload();
        }}
      >
        Reload this page
      </a>
    </div>
  );
};

export default LoadingError;
