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

import classNames from "classnames/bind";
import styles from "./filter.module.scss";

const cx = classNames.bind(styles);

export const dropdownButton = cx("dropdown-btn");
export const dropdownContentWrapper = cx("dropdown-content-wrapper");
export const dropdown = cx("dropdown-area");
export const hidden = cx("hide");
export const caretDown = cx("caret-down");

export const timePair = {
  wrapper: cx("time-pair-wrapper"),
  timeNumber: cx("time-number"),
  timeUnit: cx("time-unit"),
};

export const filterLabel = {
  top: cx("filter-label"),
  margin: cx("filter-label", "filter-label__margin-top"),
};

export const checkbox = {
  input: cx("checkbox__input"),
  label: cx("checkbox__label"),
};

export const applyBtn = {
  wrapper: cx("apply-btn__wrapper"),
  btn: cx("apply-btn__btn"),
};
