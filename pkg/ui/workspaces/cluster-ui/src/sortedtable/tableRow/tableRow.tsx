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
import classNames from "classnames/bind";
import styles from "./tableRow.module.scss";
import { ExpandableConfig, SortableColumn } from "../sortedtable";
import { RowCell } from "./rowCell";

const cx = classNames.bind(styles);

interface TableRowProps {
  columns: SortableColumn[];
  expandableConfig?: ExpandableConfig;
  firstCellBordered: boolean;
  rowClass?: (rowIndex: number) => string;
  rowIndex: number;
}

interface ExpansionProps {
  classes: string;
  rowIndex: number;
  colSpan: number;
  expandableConfig: ExpandableConfig;
}

const Expansion: React.FC<ExpansionProps> = ({
  classes,
  rowIndex,
  colSpan,
  expandableConfig,
}) => {
  const expandedAreaClasses = cx(
    "row-wrapper__row",
    "row-wrapper__row--body",
    "row-wrapper__row--expanded-area",
  );
  const cellClass = cx("row-wrapper__cell");

  return (
    // Add a zero-height empty row so that the expanded area will have the same background
    // color as the row it's expanded from, since the CSS causes row colors to alternate.
    <>
      <tr className={classes} key={"expansion" + rowIndex} />,
      <tr className={expandedAreaClasses} key={"expansionContent" + rowIndex}>
        <td />
        <td className={cellClass} colSpan={colSpan}>
          {expandableConfig.expandedContent(rowIndex)}
        </td>
      </tr>
    </>
  );
};

const ExpansionControl: React.FC<{ expanded: boolean }> = ({ expanded }) => {
  const content = expanded ? "▼" : "▶";
  const controlClass = cx(
    "row-wrapper__cell",
    "row-wrapper__cell__expansion-control",
  );
  return (
    <td className={controlClass}>
      <div>{content}</div>
    </td>
  );
};

export const TableRow: React.FC<TableRowProps> = ({
  expandableConfig,
  columns,
  firstCellBordered,
  rowIndex,
  rowClass,
}) => {
  const classes = cx(
    "row-wrapper__row",
    "row-wrapper__row--body",
    rowClass(rowIndex),
    { "body-row__row--expandable": !!expandableConfig },
  );
  const expanded = expandableConfig && expandableConfig.rowIsExpanded(rowIndex);
  const handleExpand = expandableConfig && expandableConfig.onChangeExpansion;
  const rowAction = handleExpand
    ? () => handleExpand(rowIndex, !expanded)
    : null;
  return (
    <>
      <tr className={classes} onClick={rowAction}>
        {expandableConfig && <ExpansionControl expanded={true} />}
        {columns.map((c: SortableColumn, colIndex: number) => {
          const cellClasses = cx(
            "row-wrapper__cell",
            { "cell-header": firstCellBordered && colIndex === 0 },
            c.className,
          );
          return (
            <RowCell cellClasses={cellClasses} key={"rowCell" + colIndex}>
              {c.cell(rowIndex)}
            </RowCell>
          );
        })}
      </tr>
      {expandableConfig && expandableConfig.rowIsExpanded(rowIndex) && (
        <Expansion
          expandableConfig={expandableConfig}
          rowIndex={rowIndex}
          classes={classes}
          colSpan={columns.length}
        />
      )}
    </>
  );
};
