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

import React, { ButtonHTMLAttributes } from "react";
import classNames from "classnames/bind";
import styles from "./button.module.styl";

export interface ButtonProps {
  type?: "primary" | "secondary" | "flat" | "unstyled-link";
  disabled?: boolean;
  textAlign?: "left" | "right" | "center";
  size?: "default" | "small";
  children?: React.ReactNode;
  icon?: () => React.ReactNode;
  iconPosition?: "left" | "right";
  onClick?: (event: React.MouseEvent<HTMLElement>) => void;
  className?: string;
  buttonType?: ButtonHTMLAttributes<HTMLButtonElement>["type"];
  tabIndex?: ButtonHTMLAttributes<HTMLButtonElement>["tabIndex"];
}

const cx = classNames.bind(styles);

export function Button(props: ButtonProps) {
  const {
    children,
    type,
    disabled,
    size,
    icon,
    iconPosition,
    onClick,
    className,
    buttonType,
    tabIndex,
    textAlign,
  } = props;

  const rootStyles = cx(
    "crl-button",
    `crl-button--type-${type}`,
    `crl-button--size-${size}`,
    {
      "crl-button--disabled": disabled,
    },
    className,
  );

  const renderIcon = () => {
    if (icon === undefined) {
      return null;
    }
    return (
      <div className={cx(`crl-button__icon--push-${iconPosition}`)}>
        {icon()}
      </div>
    );
  };

  return (
    <button
      onClick={onClick}
      className={rootStyles}
      disabled={disabled}
      type={buttonType}
      tabIndex={tabIndex}
    >
      <div className={cx("crl-button__container")}>
        {iconPosition === "left" && renderIcon()}
        <div className={cx("crl-button__content")} style={{ textAlign }}>
          {children}
        </div>
        {iconPosition === "right" && renderIcon()}
      </div>
    </button>
  );
}

Button.defaultProps = {
  onClick: () => {},
  type: "primary",
  disabled: false,
  size: "default",
  className: "",
  iconPosition: "left",
  buttonType: "button",
  textAlign: "left",
};
