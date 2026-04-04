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
import classnames from "classnames/bind";

import { OutsideEventHandler } from "../outsideEventHandler";
import styles from "./dropdown.module.scss";
import { Button, ButtonProps } from "src/button";
import { CaretDown } from "@cockroachlabs/icons";

const cx = classnames.bind(styles);

export interface DropdownOption<T = string> {
  value: T;
  name: React.ReactNode | string;
  disabled?: boolean;
}

export interface DropdownItemProps<T> {
  children: React.ReactNode;
  value: T;
  onClick: (value: T) => void;
  disabled?: boolean;
  className?: string;
}
export interface DropdownProps<T> {
  items: Array<DropdownOption<T>>;
  onChange: (item: DropdownOption<T>["value"]) => void;
  children?: React.ReactNode;
  customToggleButton?: React.ReactChild;
  customToggleButtonOptions?: Partial<ButtonProps>;
  menuPosition?: "left" | "right";
  className?: string;
  itemsClassname?: string;
}

interface DropdownState {
  isOpen: boolean;
}

interface DropdownButtonProps {
  children: React.ReactNode;
  isOpen: boolean;
  onClick?: (event: React.MouseEvent<HTMLElement>) => void;
  customProps?: Partial<ButtonProps>;
}

export const DropdownButton: React.FC<DropdownButtonProps> = ({
  children,
  customProps = {},
}) => {
  return (
    <Button
      type="secondary"
      size="default"
      iconPosition="right"
      icon={<CaretDown />}
      {...customProps}
    >
      {children}
    </Button>
  );
};

function DropdownItem<T = string>(props: DropdownItemProps<T>) {
  const { children, value, onClick, disabled, className } = props;
  return (
    <div
      onClick={() => onClick(value)}
      className={cx(
        "crl-dropdown__item",
        {
          "crl-dropdown__item--disabled": disabled,
        },
        className,
      )}
    >
      {children}
    </div>
  );
}

export class Dropdown<T = string> extends React.Component<
  DropdownProps<T>,
  DropdownState
> {
  state = {
    isOpen: false,
  };

  handleMenuOpen = (): void => {
    this.setState({
      isOpen: !this.state.isOpen,
    });
  };

  changeMenuState = (nextState: boolean): void => {
    this.setState({
      isOpen: nextState,
    });
  };

  handleItemSelection = (value: T): void => {
    this.props.onChange(value);
    this.handleMenuOpen();
  };

  renderDropdownToggleButton = (): React.ReactChild => {
    const {
      children,
      customToggleButton,
      customToggleButtonOptions,
    } = this.props;
    const { isOpen } = this.state;

    if (customToggleButton) {
      return customToggleButton;
    } else {
      return (
        <DropdownButton isOpen={isOpen} customProps={customToggleButtonOptions}>
          {children}
        </DropdownButton>
      );
    }
  };

  render(): React.ReactElement {
    const {
      items,
      menuPosition = "left",
      className,
      itemsClassname,
    } = this.props;
    const { isOpen } = this.state;

    const menuStyles = cx(
      "crl-dropdown__menu",
      `crl-dropdown__menu--align-${menuPosition}`,
      {
        "crl-dropdown__menu--open": isOpen,
      },
    );

    const menuItems = items.map((menuItem, idx) => (
      <DropdownItem<T>
        value={menuItem.value}
        onClick={this.handleItemSelection}
        key={idx}
        disabled={menuItem.disabled}
        className={itemsClassname}
      >
        {menuItem.name}
      </DropdownItem>
    ));

    return (
      <div className={cx("crl-dropdown", className)}>
        <OutsideEventHandler onOutsideClick={() => this.changeMenuState(false)}>
          <div
            className={cx("crl-dropdown__handler")}
            onClick={this.handleMenuOpen}
          >
            {this.renderDropdownToggleButton()}
          </div>
          <div className={cx("crl-dropdown__overlay")}>
            <div className={menuStyles}>
              <div className={cx("crl-dropdown__container")}>{menuItems}</div>
            </div>
          </div>
        </OutsideEventHandler>
      </div>
    );
  }
}
