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

import EyeOff from "assets/eye-off.svg";
import Eye from "assets/eye.svg";
import cn from "classnames";
import React from "react";
import { Button } from "../button";
import "./input.styl";

interface PasswordInputProps {
  onChange: (value: string) => void;
  value: string;
  placeholder?: string;
  className?: string;
  name?: string;
  label?: string;
}

interface PasswordInputState {
  showPassword?: boolean;
}

export class PasswordInput extends React.Component<
  PasswordInputProps,
  PasswordInputState
> {
  state = {
    showPassword: false,
  };

  handleOnTextChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    this.props.onChange(value);
  };

  togglePassword = () => {
    this.setState({
      showPassword: !this.state.showPassword,
    });
  };

  renderPasswordIcon = (showPassword: boolean) => (
    <Button
      tabIndex={-1}
      type="flat"
      onClick={this.togglePassword}
      className="crl-button__show-password"
    >
      <img src={showPassword ? EyeOff : Eye} alt="Toggle Password" />
    </Button>
  );

  render() {
    const { placeholder, className, name, label, value } = this.props;
    const { showPassword } = this.state;
    const inputType = showPassword ? "text" : "password";

    const classes = cn(className, "crl-input", "crl-input__password");
    return (
      <div className="crl-input__wrapper">
        {label && (
          <label htmlFor={name} className="crl-input__label">
            {label}
          </label>
        )}
        <input
          name={name}
          type={inputType}
          value={value}
          placeholder={placeholder}
          className={classes}
          onChange={this.handleOnTextChange}
        />
        {this.renderPasswordIcon(showPassword)}
      </div>
    );
  }
}
