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
import Select, { Props, OptionsType } from "react-select";
import { noop } from "lodash";
import { CheckboxOption } from "../multiSelectCheckbox/multiSelectCheckbox";
import { filterLabel } from "../queryFilter/filterClasses";
import { StylesConfig } from "react-select/src/styles";

export type FilterCheckboxOptionItem = { label: string; value: string };
export type FilterCheckboxOptionsType = OptionsType<FilterCheckboxOptionItem>;

export type FilterCheckboxOptionProps = {
  label: string;
  // onSelectionChanged callback function is called with all selected options.
  onSelectionChanged?: (options: OptionsType<FilterCheckboxOptionItem>) => void;
  triggerClear?: (fn: () => void) => void;
} & Props;

export const FilterCheckboxOption = (props: FilterCheckboxOptionProps) => {
  const {
    label,
    onSelectionChanged = noop,
    options,
    placeholder,
    ...selectProps
  } = props;

  const customStyles: StylesConfig<FilterCheckboxOptionItem, true> = {
    container: provided => ({
      ...provided,
      border: "none",
    }),
    option: (provided, state) => ({
      ...provided,
      backgroundColor: state.isSelected ? "#DEEBFF" : provided.backgroundColor,
      color: "#394455",
    }),
    control: provided => ({
      ...provided,
      width: "100%",
      borderColor: "#C0C6D9",
    }),
    dropdownIndicator: provided => ({
      ...provided,
      color: "#C0C6D9",
    }),
    singleValue: provided => ({
      ...provided,
      color: "#475872",
    }),
    menu: provided => ({
      ...provided,
      zIndex: 3,
    }),
  };

  return (
    <div>
      <div className={filterLabel.margin}>{label}</div>
      <Select
        {...selectProps}
        isMulti
        options={options}
        placeholder={placeholder}
        onChange={onSelectionChanged}
        hideSelectedOptions={false}
        closeMenuOnSelect={false}
        components={{ Option: CheckboxOption }}
        styles={customStyles}
      />
    </div>
  );
};
