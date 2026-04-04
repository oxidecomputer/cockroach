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
import styles from "./anchor.module.scss";

interface AnchorProps {
  onClick?: () => void;
  href?: string;
  target?: "_blank" | "_parent" | "_self";
  className?: string;
}

const cx = classnames.bind(styles);

export function Anchor(props: React.PropsWithChildren<AnchorProps>) {
  const { href, target, children, onClick, className } = props;
  return (
    <a
      className={cx("crl-anchor", className)}
      href={href}
      target={target}
      onClick={onClick}
    >
      {children}
    </a>
  );
}

Anchor.defaultProps = {
  target: "_blank",
  className: "",
};
