// Copyright 2022 The Cockroach Authors.
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
import { DropdownButton } from "../dropdown";
import { OutsideEventHandler } from "../outsideEventHandler";
import classnames from "classnames/bind";
import styles from "../dropdown/dropdown.module.scss";
import { applyBtn } from "../queryFilter/filterClasses";
import { Button } from "../button";

const cx = classnames.bind(styles);

type FilterDropdownProps = React.PropsWithChildren<{
  className?: string;
  label: string;
  onSubmit: () => void;
}>;

export const FilterDropdown = ({
  className,
  label,
  onSubmit,
  children,
}: FilterDropdownProps) => {
  const [isOpen, setIsOpen] = React.useState<boolean>(false);
  const toggleMenuState = React.useCallback(() => {
    setIsOpen(!isOpen);
  }, [isOpen]);

  const onSubmitCallback = React.useCallback(() => {
    onSubmit();
    setIsOpen(false);
  }, [onSubmit]);

  const menuStyles = cx(
    "crl-dropdown__menu",
    `crl-dropdown__menu--align-left`,
    {
      "crl-dropdown__menu--open": isOpen,
    },
  );

  return (
    <div
      className={cx("crl-dropdown", className)}
      onClick={event => event.stopPropagation()}
    >
      <OutsideEventHandler onOutsideClick={() => setIsOpen(false)}>
        <div className={cx("crl-dropdown__handler")} onClick={toggleMenuState}>
          <DropdownButton isOpen={true}>{label}</DropdownButton>
        </div>
        <div className={cx("crl-dropdown__overlay")}>
          <div className={menuStyles}>
            <div
              className={cx(
                "crl-dropdown__container",
                "crl-dropdown__container__wrapped",
              )}
            >
              {children}
              <div className={applyBtn.wrapper}>
                <Button
                  className={applyBtn.btn}
                  textAlign="center"
                  onClick={onSubmitCallback}
                >
                  Apply
                </Button>
              </div>
            </div>
          </div>
        </div>
      </OutsideEventHandler>
    </div>
  );
};
