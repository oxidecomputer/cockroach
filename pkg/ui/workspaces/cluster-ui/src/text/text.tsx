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
import classNames from "classnames/bind";

import styles from "./text.module.scss";

export interface TextProps {
  textType?: TextTypes;
  disabled?: boolean;
  children: React.ReactNode;
  className?: string;
  noWrap?: boolean;
}

export enum TextTypes {
  Heading1,
  Heading2,
  Heading3,
  Heading4,
  Heading5,
  Heading6,
  Body,
  BodyStrong,
  Caption,
  CaptionStrong,
  Code,
}

const getClassByTextType = (textType: TextTypes) => {
  switch (textType) {
    case TextTypes.Heading1:
      return "text--heading-1";
    case TextTypes.Heading2:
      return "text--heading-2";
    case TextTypes.Heading3:
      return "text--heading-3";
    case TextTypes.Heading4:
      return "text--heading-4";
    case TextTypes.Heading5:
      return "text--heading-5";
    case TextTypes.Heading6:
      return "text--heading-6";
    case TextTypes.Body:
      return "text--body";
    case TextTypes.BodyStrong:
      return "text--body-strong";
    case TextTypes.Caption:
      return "text--caption";
    case TextTypes.CaptionStrong:
      return "text--caption-strong";
    case TextTypes.Code:
      return "text--code";
    default:
      return "text--body";
  }
};

const cx = classNames.bind(styles);

const getElementByTextType = (textType: TextTypes) => {
  switch (textType) {
    case TextTypes.Heading1:
      return "h1";
    case TextTypes.Heading2:
      return "h2";
    case TextTypes.Heading3:
      return "h3";
    case TextTypes.Heading4:
      return "h4";
    case TextTypes.Heading5:
      return "h5";
    case TextTypes.Heading6:
      return "h6";
    case TextTypes.Body:
    case TextTypes.BodyStrong:
    case TextTypes.Caption:
    case TextTypes.CaptionStrong:
    case TextTypes.Code:
    default:
      return "span";
  }
};

export function Text(props: TextProps) {
  const { textType, disabled, noWrap, className } = props;
  const textTypeClass = cx(
    "text",
    getClassByTextType(textType),
    {
      "text--disabled": disabled,
      "text--no-wrap": noWrap,
    },
    className,
  );
  const elementName = getElementByTextType(textType);

  const componentProps = {
    className: textTypeClass,
  };

  return React.createElement(elementName, componentProps, props.children);
}

Text.defaultProps = {
  textType: TextTypes.Body,
  disabled: false,
  className: "",
  noWrap: false,
};
