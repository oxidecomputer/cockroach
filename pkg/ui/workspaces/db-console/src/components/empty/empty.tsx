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

import heroBannerLp from "assets/heroBannerLp.png";
import React from "react";
import classnames from "classnames/bind";
import { Anchor, Button, Text, TextTypes } from "src/components";
import styles from "./empty.module.styl";

const cx = classnames.bind(styles);

interface IMainEmptyProps {
  title?: string;
  description?: string;
  label?: React.ReactNode;
  link?: string;
  anchor?: string;
  backgroundImage?: string;
}

type OnClickXORHref =
  | {
      onClick?: () => void;
      buttonHref?: never;
    }
  | {
      onClick?: never;
      buttonHref?: string;
    };

export type EmptyProps = OnClickXORHref & IMainEmptyProps;

export const Empty = ({
  title,
  description,
  anchor,
  label,
  link,
  backgroundImage,
  onClick,
  buttonHref,
}: EmptyProps) => (
  <div
    className={cx("cl-empty-view")}
    style={{ backgroundImage: `url(${backgroundImage})` }}
  >
    <Text className={cx("cl-empty-view__title")} textType={TextTypes.Heading3}>
      {title}
    </Text>
    <div className={cx("cl-empty-view__content")}>
      <main className={cx("cl-empty-view__main")}>
        <Text
          textType={TextTypes.Body}
          className={cx("cl-empty-view__main--text")}
        >
          {description}
          {link && (
            <Anchor href={link} className={cx("cl-empty-view__main--anchor")}>
              {anchor}
            </Anchor>
          )}
        </Text>
      </main>
      <footer className={cx("cl-empty-view__footer")}>
        <Button
          type="primary"
          onClick={() =>
            buttonHref ? window.open(buttonHref) : onClick && onClick()
          }
        >
          {label}
        </Button>
      </footer>
    </div>
  </div>
);

Empty.defaultProps = {
  backgroundImage: heroBannerLp,
  anchor: "Learn more",
  label: "Learn more",
  title: "No results",
};
