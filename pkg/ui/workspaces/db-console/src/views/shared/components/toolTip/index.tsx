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

import { Tooltip } from "antd";
import React from "react";
import { AbstractTooltipProps } from "antd/es/tooltip";
import classNames from "classnames/bind";

import styles from "./tooltip.module.styl";

interface ToolTipWrapperProps extends AbstractTooltipProps {
  text: React.ReactNode;
  short?: boolean;
  children?: React.ReactNode;
}

const cx = classNames.bind(styles);

/**
 * ToolTipWrapper wraps its children with an area that detects mouseover events
 * and, when hovered, displays a floating tooltip to the immediate right of
 * the wrapped element.
 *
 * Note that the child element itself must be wrappable; certain CSS attributes
 * such as "float" will render parent elements unable to properly wrap their
 * contents.
 */

export const ToolTipWrapper = (props: ToolTipWrapperProps) => {
  const { text, children, placement = "bottom" } = props;
  const overlayClassName = cx("tooltip-wrapper", "tooltip__preset--white");
  return (
    <Tooltip
      title={text}
      placement={placement}
      overlayClassName={overlayClassName}
      {...props}
    >
      {children}
    </Tooltip>
  );
};

ToolTipWrapper.defaultProps = {
  placement: "bottom",
};
