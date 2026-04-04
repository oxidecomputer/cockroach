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
import { Redirect, match as Match } from "react-router-dom";
import { StatementLinkTarget } from "@cockroachlabs/cluster-ui";
import { getMatchParamByName } from "src/util/query";
import {
  appAttr,
  databaseAttr,
  implicitTxnAttr,
  statementAttr,
} from "src/util/constants";

type Props = {
  match: Match;
};

// RedirectToStatementDetails is designed to route old versions of StatementDetails routes
// where app and database are route params, to the new StatementDetails route.
export function RedirectToStatementDetails({ match }: Props) {
  const linkProps = {
    statementFingerprintID: getMatchParamByName(match, statementAttr),
    app: getMatchParamByName(match, appAttr),
    implicitTxn: getMatchParamByName(match, implicitTxnAttr) === "true",
    database: getMatchParamByName(match, databaseAttr),
  };

  return <Redirect to={StatementLinkTarget(linkProps)} />;
}
