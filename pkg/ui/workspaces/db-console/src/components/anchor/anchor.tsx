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
import classnames from "classnames/bind";
import styles from "./anchor.module.styl";

type AnchorProps = React.DetailedHTMLProps<
  React.AnchorHTMLAttributes<HTMLAnchorElement>,
  HTMLAnchorElement
>;

const cx = classnames.bind(styles);

export function Anchor({
  target = "_blank",
  className,
  children,
  ...props
}: AnchorProps) {
  return (
    <a {...props} className={cx("crl-anchor", className)} target={target}>
      {children}
    </a>
  );
}
