// Copyright 2019 The Cockroach Authors.
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

import * as React from "react";
import cn from "classnames";

import { Text, TextTypes } from "src/components";

import "./sideNavigation.styl";

export interface SideNavigationProps {
  children: React.ReactNode;
  className?: string;
}

export interface NavigationItem {
  disabled?: boolean;
  isActive?: boolean;
  className?: string;
  children: React.ReactNode;
}

export function NavigationItem(props: NavigationItem) {
  const { children, isActive, disabled } = props;
  let textType = TextTypes.Body;

  if (isActive) {
    textType = TextTypes.BodyStrong;
  }

  const classes = cn("side-navigation__navigation-item", {
    "side-navigation__navigation-item--active": isActive,
    "side-navigation__navigation-item--disabled": disabled,
  });

  return (
    <li className={classes}>
      <Text textType={textType}>{children}</Text>
    </li>
  );
}

NavigationItem.defaultProps = {
  isActive: false,
  disabled: false,
};

SideNavigation.Item = NavigationItem;

export function SideNavigation(props: SideNavigationProps) {
  return (
    <nav className="side-navigation">
      <ul className="side-navigation--list">{props.children}</ul>
    </nav>
  );
}
