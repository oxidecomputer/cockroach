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

import React, { useContext } from "react";
import classnames from "classnames/bind";
import styles from "./pageConfig.module.scss";
import { CockroachCloudContext } from "../contexts";

export interface PageConfigProps {
  layout?: "list" | "spread";
  children?: React.ReactNode;
  className?: string;
}

const cx = classnames.bind(styles);

export function PageConfig(props: PageConfigProps): React.ReactElement {
  const isCockroachCloud = useContext(CockroachCloudContext);

  const classes = cx({
    "page-config__list": props.layout !== "spread",
    "page-config__spread": props.layout === "spread",
  });

  return (
    <div
      className={`${cx("page-config", {
        "page-config__white-background": isCockroachCloud,
      })} ${props.className}`}
    >
      <ul className={classes}>{props.children}</ul>
    </div>
  );
}

export interface PageConfigItemProps {
  children?: React.ReactNode;
  className?: string;
}

export function PageConfigItem(props: PageConfigItemProps): React.ReactElement {
  return (
    <li className={`${cx("page-config__item")} ${props.className}`}>
      {props.children}
    </li>
  );
}
