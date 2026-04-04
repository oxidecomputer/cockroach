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
import { EmptyTable, EmptyTableProps } from "src/empty";
import { Anchor } from "src/anchor";
import { statementsTable } from "src/util";
import magnifyingGlassImg from "../assets/emptyState/magnifying-glass.svg";
import emptyTableResultsImg from "../assets/emptyState/empty-table-results.svg";

export const EmptyStatementsPlaceholder: React.FC<{
  isEmptySearchResults: boolean;
}> = ({ isEmptySearchResults }) => {
  const footer = (
    <Anchor href={statementsTable} target="_blank">
      Learn more about statements
    </Anchor>
  );

  const emptyPlaceholderProps: EmptyTableProps = isEmptySearchResults
    ? {
        title:
          "No SQL statements match your search in the selected time interval",
        icon: magnifyingGlassImg,
        footer,
      }
    : {
        title: "No SQL statements in the selected time interval",
        icon: emptyTableResultsImg,
        footer,
      };
  return <EmptyTable {...emptyPlaceholderProps} />;
};
