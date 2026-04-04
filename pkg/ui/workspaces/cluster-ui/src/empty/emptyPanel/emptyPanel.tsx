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
import { Heading, Text, Button } from "@cockroachlabs/ui-components";
import styles from "./emptyPanel.module.scss";
import { Anchor } from "../../anchor";
import heroBannerLp from "../../assets/heroBannerLp.png";

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

export type EmptyPanelProps = OnClickXORHref & IMainEmptyProps;

export const EmptyPanel: React.FC<EmptyPanelProps> = ({
  title = "No results",
  description,
  anchor = "Learn more",
  label = "Learn more",
  link,
  backgroundImage = heroBannerLp,
  onClick,
  buttonHref,
}) => (
  <div
    className={cx("cl-empty-view")}
    style={{ backgroundImage: `url(${backgroundImage})` }}
  >
    <Heading type="h2">{title}</Heading>
    <div className={cx("cl-empty-view__content")}>
      <main className={cx("cl-empty-view__main")}>
        <Text>
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
          intent="primary"
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
