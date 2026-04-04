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

import React, { FunctionComponent, ReactElement } from "react";
import { Link } from "react-router-dom";
import classnames from "classnames/bind";
import styles from "./breadcrumbs.module.scss";

export interface BreadcrumbItem {
  name: string;
  link: string;
  onClick?: () => void;
}

interface BreadcrumbsProps {
  items: BreadcrumbItem[];
  divider?: ReactElement;
}

const cx = classnames.bind(styles);

export const Breadcrumbs: FunctionComponent<BreadcrumbsProps> = ({
  items,
  divider = "/",
}) => {
  if (items.length === 0) {
    return null;
  }
  const lastItem = items.pop();
  return (
    <div className={cx("breadcrumbs")}>
      {items.map(({ link, name, onClick = () => {} }) => (
        <div className={cx("breadcrumbs__item")} key={link}>
          <Link
            className={cx("breadcrumbs__item--link")}
            to={link}
            onClick={onClick}
          >
            {name}
          </Link>
          <span className={cx("breadcrumbs__item--divider")}>{divider}</span>
        </div>
      ))}
      <div className={cx("breadcrumbs__item", "breadcrumbs__item--last")}>
        {lastItem?.name}
      </div>
    </div>
  );
};
