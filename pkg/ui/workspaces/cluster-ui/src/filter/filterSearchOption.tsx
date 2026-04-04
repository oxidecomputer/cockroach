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
import { Search } from "../search";
import { filterLabel } from "../queryFilter/filterClasses";

export type FilterSearchOptionProps = {
  label: string;
  onChanged?: (value: string) => void;
  value?: string;
};

export const FilterSearchOption = (props: FilterSearchOptionProps) => {
  const { label, onChanged, value } = props;
  return (
    <div>
      <div className={filterLabel.margin}>{label}</div>
      <Search
        onChange={onChanged}
        renderSuffix={false}
        placeholder="Search"
        value={value}
      />
    </div>
  );
};
