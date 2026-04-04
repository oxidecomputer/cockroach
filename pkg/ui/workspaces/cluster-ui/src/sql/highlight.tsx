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

import hljs from "highlight.js/lib/core";
import sqlLangSyntax from "highlight.js/lib/languages/pgsql";
import React from "react";
import classNames from "classnames/bind";
import styles from "./sqlhighlight.module.scss";
import { SqlBoxProps } from "./box";

const cx = classNames.bind(styles);

hljs.registerLanguage("sql", sqlLangSyntax);
hljs.configure({
  tabReplace: "  ",
});

export class Highlight extends React.Component<SqlBoxProps> {
  preNode: React.RefObject<HTMLPreElement> = React.createRef();

  shouldComponentUpdate(newProps: SqlBoxProps) {
    return newProps.value !== this.props.value;
  }

  componentDidMount() {
    hljs.highlightBlock(this.preNode.current);
  }

  componentDidUpdate() {
    hljs.highlightBlock(this.preNode.current);
  }

  renderZone = () => {
    const { zone } = this.props;
    const zoneConfig = zone.zone_config;
    return (
      <span className={cx("sql-highlight", "hljs")}>
        <span className="hljs-keyword">CONFIGURE ZONE USING</span>
        <br />
        <span className="hljs-label">range_min_bytes = </span>
        <span className="hljs-built_in">{`${String(
          zoneConfig.range_min_bytes,
        )},`}</span>
        <br />
        <span className="hljs-label">range_max_bytes = </span>
        <span className="hljs-built_in">{`${String(
          zoneConfig.range_max_bytes,
        )},`}</span>
        <br />
        <span className="hljs-label">gc.ttlseconds = </span>
        <span className="hljs-built_in">{`${zoneConfig.gc.ttl_seconds},`}</span>
        <br />
        <span className="hljs-label">num_replicas = </span>
        <span className="hljs-built_in">{`${zoneConfig.num_replicas},`}</span>
        <br />
        <span className="hljs-label">constraints = [&apos;</span>
        <span className="hljs-built_in">{String(zoneConfig.constraints)}</span>
        &apos;],
        <br />
        <span className="hljs-label">lease_preferences = [[&apos;</span>
        <span className="hljs-built_in">
          {String(zoneConfig.lease_preferences)}
        </span>
        &apos;]]
      </span>
    );
  };

  render() {
    const { value, zone } = this.props;
    return (
      <>
        <span className={cx("sql-highlight")} ref={this.preNode}>
          {value}
        </span>
        {zone && (
          <>
            <div className={cx("highlight-divider")} />
            {this.renderZone()}
          </>
        )}
      </>
    );
  }
}
