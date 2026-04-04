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

import { Button, Form, Input } from "antd";
import { InputProps } from "antd/lib/input";
import CancelIcon from "assets/cancel.svg";
import SearchIcon from "assets/search.svg";
import React from "react";
import classNames from "classnames/bind";
import styles from "./search.module.styl";

interface ISearchProps {
  onSubmit: (value: string) => void;
  onClear?: () => void;
  defaultValue?: string;
}

interface ISearchState {
  value: string;
  submitted: boolean;
  submit?: boolean;
}

type TSearchProps = ISearchProps & InputProps;

const cx = classNames.bind(styles);

export class Search extends React.Component<TSearchProps, ISearchState> {
  state = {
    value: this.props.defaultValue || "",
    submitted: false,
  };

  onSubmit = () => {
    const { value } = this.state;
    const { onSubmit } = this.props;
    onSubmit(value);
    if (value.length > 0) {
      this.setState({
        submitted: true,
      });
    }
  };

  onChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value: string = event.target.value;
    const submitted = value.length === 0;
    this.setState({ value, submitted }, () => submitted && this.onSubmit());
  };

  onClear = () => {
    const { onClear } = this.props;
    this.setState({ value: "", submit: false });
    onClear();
  };

  renderSuffix = () => {
    const { value, submitted } = this.state;
    if (value.length > 0) {
      if (submitted) {
        return (
          <Button
            onClick={this.onClear}
            type="default"
            className={cx("_clear-search")}
          >
            <img className={cx("_suffix-icon")} src={CancelIcon} />
          </Button>
        );
      }
      return (
        <Button
          type="default"
          htmlType="submit"
          className={cx("_submit-search")}
        >
          Enter
        </Button>
      );
    }
    return null;
  };

  render() {
    const { value, submitted } = this.state;
    const className = submitted ? cx("_submitted") : "";

    /*
      current antdesign (3.25.3) has Input.d.ts incompatible with latest @types/react
      temporary fix, remove as soon antd can be updated
    */

    const MyInput = Input as any;

    return (
      <Form onSubmit={this.onSubmit} className={cx("_search-form")}>
        <Form.Item>
          <MyInput
            className={className}
            placeholder="Search Statement"
            onChange={this.onChange}
            prefix={<img className={cx("_prefix-icon")} src={SearchIcon} />}
            suffix={this.renderSuffix()}
            value={value}
            {...this.props}
          />
        </Form.Item>
      </Form>
    );
  }
}
