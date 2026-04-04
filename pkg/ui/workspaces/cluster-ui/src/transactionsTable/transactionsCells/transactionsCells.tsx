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
import { Link } from "react-router-dom";
import { getHighlightedText } from "src/highlightedText";
import { Tooltip } from "@cockroachlabs/ui-components";
import { limitText } from "../utils";
import classNames from "classnames/bind";
import statementsStyles from "../../statementsTable/statementsTableContent.module.scss";
import transactionsCellsStyles from "./transactionsCells.module.scss";
import { TransactionLinkTarget } from "../transactionsTable";

const statementsCx = classNames.bind(statementsStyles);
const ownCellStyles = classNames.bind(transactionsCellsStyles);
const descriptionClassName = statementsCx("cl-table-link__description");
const textWrapper = ownCellStyles("text-wrapper");
const hoverAreaClassName = ownCellStyles("hover-area");

interface TextCellProps {
  transactionText: string;
  transactionSummary: string;
  aggregatedTs: string;
  transactionFingerprintId: string;
  search: string;
}

export const transactionLink = ({
  transactionText,
  transactionSummary,
  aggregatedTs,
  transactionFingerprintId,
  search,
}: TextCellProps): React.ReactElement => {
  const linkProps = {
    aggregatedTs,
    transactionFingerprintId,
  };

  return (
    <Link to={TransactionLinkTarget(linkProps)}>
      <div>
        <Tooltip
          placement="bottom"
          content={
            <pre className={descriptionClassName}>
              {getHighlightedText(
                transactionText,
                search,
                true /* hasDarkBkg */,
              )}
            </pre>
          }
        >
          <div className={textWrapper}>
            <div className={hoverAreaClassName}>
              {getHighlightedText(
                limitText(transactionSummary, 200),
                search,
                false /* hasDarkBkg */,
                true /* isOriginalText */,
              )}
            </div>
          </div>
        </Tooltip>
      </div>
    </Link>
  );
};
