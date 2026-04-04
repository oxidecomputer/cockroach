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

export interface PaginationSettings {
  pageSize?: number;
  current: number;
  total?: number;
}

export interface ResultsPerPageLabelProps {
  pagination: PaginationSettings;
  pageName: string;
  selectedApp?: string;
  search?: string;
}

export const ResultsPerPageLabel: React.FC<ResultsPerPageLabelProps> = ({
  pagination: { pageSize, current, total },
  pageName,
  selectedApp = "",
  search,
}) => {
  const getPageStart = (pageSize: number, current: number) =>
    pageSize * current;
  const start = Math.max(
    getPageStart(pageSize, current > 0 ? current - 1 : current),
    0,
  );
  const recountedStart = total > 0 ? start + 1 : start;
  const end = Math.min(getPageStart(pageSize, current), total);
  const label =
    (search && search.length > 0) || selectedApp.length > 0
      ? "results"
      : pageName;
  if (end === 0) {
    return <>{`0 ${label}`}</>;
  }
  return <>{`${recountedStart}-${end} of ${total} ${label}`}</>;
};
