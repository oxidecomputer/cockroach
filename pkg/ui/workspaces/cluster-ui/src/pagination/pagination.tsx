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
import { Pagination as AntPagination } from "antd";
import { PaginationProps as AntPaginationProps } from "antd/lib/pagination";
import classNames from "classnames/bind";
import styles from "./pagination.module.scss";

const cx = classNames.bind(styles);

export type PaginationProps = Pick<
  AntPaginationProps,
  "pageSize" | "current" | "total" | "onChange"
>;

export const Pagination: React.FC<PaginationProps> = props => {
  const itemRenderer = React.useCallback(
    (
      _page: number,
      type: "page" | "prev" | "next" | "jump-prev" | "jump-next",
      originalElement: React.ReactNode,
    ) => {
      switch (type) {
        case "jump-prev":
        case "jump-next":
          return (
            <div className={cx("_pg-jump")}>
              <span className={cx("_jump-dots")}>•••</span>
            </div>
          );
        default:
          return originalElement;
      }
    },
    [],
  );

  return (
    <AntPagination
      {...props}
      size="small"
      itemRender={itemRenderer}
      hideOnSinglePage
      className={cx("root")}
    />
  );
};
