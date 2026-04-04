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

import * as hljs from "highlight.js";
import React from "react";
import classNames from "classnames/bind";
import styles from "./sqlhighlight.module.styl";
import { SqlBoxProps } from "./box";

const cx = classNames.bind(styles);

export class Highlight extends React.Component<SqlBoxProps> {
  preNode: React.RefObject<HTMLPreElement> = React.createRef();
  preNodeSecondary: React.RefObject<HTMLPreElement> = React.createRef();

  shouldComponentUpdate(newProps: SqlBoxProps) {
    return newProps.value !== this.props.value;
  }

  componentDidMount() {
    hljs.configure({
      tabReplace: "  ",
      languages: ["sql"],
    });
    hljs.highlightBlock(this.preNode.current);
    if (this.preNodeSecondary.current) {
      hljs.highlightBlock(this.preNodeSecondary.current);
    }
  }

  componentDidUpdate() {
    hljs.highlightBlock(this.preNode.current);
    if (this.preNodeSecondary.current) {
      hljs.highlightBlock(this.preNodeSecondary.current);
    }
  }

  render() {
    const { value, secondaryValue } = this.props;
    return (
      <>
        <span className={cx("sql-highlight")} ref={this.preNode}>
          {value}
        </span>
        {secondaryValue && (
          <>
            <div className={cx("highlight-divider")} />
            <span className={cx("sql-highlight")} ref={this.preNodeSecondary}>
              {secondaryValue}
            </span>
          </>
        )}
      </>
    );
  }
}
