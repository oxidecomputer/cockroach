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
import { Drawer, Button, Divider } from "antd";
import { Link } from "react-router-dom";
import classNames from "classnames/bind";
import styles from "./drawer.module.styl";

const cx = classNames.bind(styles);

interface IDrawerProps {
  visible: boolean;
  onClose: () => void;
  details?: boolean;
  children?: React.ReactNode | string;
  data: any;
}

const openDetails = (data: any) => {
  const base =
    data.app && data.app.length > 0
      ? `/statements/${data.app}/${data.implicitTxn}`
      : `/statement/${data.implicitTxn}`;
  const link = `${base}/${encodeURIComponent(data.statement)}`;
  return <Link to={link}>View statement details</Link>;
};

export const DrawerComponent = ({
  visible,
  onClose,
  children,
  data,
  details,
  ...props
}: IDrawerProps) => (
  <Drawer
    title={
      <div className={cx("__actions")}>
        <Button type="default" ghost block onClick={onClose}>
          Close
        </Button>
        {details && (
          <React.Fragment>
            <Divider type="vertical" />
            {openDetails(data)}
          </React.Fragment>
        )}
      </div>
    }
    placement="bottom"
    closable={false}
    onClose={onClose}
    visible={visible}
    className={cx("drawer--preset-black")}
    // getContainer={false}
    {...props}
  >
    {children}
  </Drawer>
);
