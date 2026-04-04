// Copyright 2018 The Cockroach Authors.
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
import { Highlight } from "./highlight";
import classNames from "classnames/bind";

import styles from "./sqlhighlight.module.scss";
import * as protos from "@cockroachlabs/crdb-protobuf-client";
import { FormatQuery } from "src/util";

export enum SqlBoxSize {
  small = "small",
  large = "large",
}

export interface SqlBoxProps {
  value: string;
  zone?: protos.cockroach.server.serverpb.DatabaseDetailsResponse;
  className?: string;
  size?: SqlBoxSize;
  format?: boolean;
}

const cx = classNames.bind(styles);

export class SqlBox extends React.Component<SqlBoxProps> {
  preNode: React.RefObject<HTMLPreElement> = React.createRef();
  render(): React.ReactElement {
    const value = this.props.format
      ? FormatQuery(this.props.value)
      : this.props.value;
    const sizeClass = this.props.size ? this.props.size : "";
    return (
      <div className={cx("box-highlight", this.props.className, sizeClass)}>
        <Highlight {...this.props} value={value} />
      </div>
    );
  }
}
