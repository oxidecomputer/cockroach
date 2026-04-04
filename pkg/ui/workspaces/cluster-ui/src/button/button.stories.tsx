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

/* eslint-disable react/jsx-key */
import React from "react";
import { storiesOf } from "@storybook/react";

import { Button, ButtonProps } from "src/button";
import { CaretDown } from "@cockroachlabs/icons";
import { Text, TextTypes } from "../text";

const sizes: ButtonProps["size"][] = ["default", "small"];
const types: ButtonProps["type"][] = [
  "primary",
  "secondary",
  "flat",
  "unstyled-link",
];
const icons: ButtonProps["icon"][] = [<CaretDown />, undefined];
const iconPositions: ButtonProps["iconPosition"][] = ["right", "left"];

storiesOf("Button", module)
  .addDecorator(renderChild => (
    <div style={{ padding: "12px", display: "flex" }}>{renderChild()}</div>
  ))
  .add("default", () => <Button>Caption</Button>)
  .add("examples", () => {
    const buttons = types.map(buttonType => {
      const buttonsPerSize = sizes.map(size => {
        const items = icons
          .map(buttonIcon => {
            return iconPositions.map(iconPosition => {
              return (
                <Button
                  type={buttonType}
                  size={size}
                  icon={buttonIcon}
                  iconPosition={iconPosition}
                >
                  Sample text
                </Button>
              );
            });
          })
          .reduce((ac, el) => [...ac, ...el], []);

        return (
          <div
            style={{
              display: "flex",
              flexDirection: "row",
              justifyContent: "space-around",
              margin: "24px 0",
            }}
          >
            {React.Children.toArray(items)}
          </div>
        );
      });

      return (
        <div>
          <Text textType={TextTypes.Heading3}>{buttonType} type</Text>
          <div style={{ display: "flex", flexDirection: "column" }}>
            {React.Children.toArray(buttonsPerSize)}
          </div>
        </div>
      );
    });

    return (
      <div style={{ display: "flex", flexDirection: "column", width: "100%" }}>
        {React.Children.toArray(buttons)}
      </div>
    );
  });
