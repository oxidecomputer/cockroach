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
import { storiesOf } from "@storybook/react";
import { MemoryRouter } from "react-router-dom";
import { noop } from "lodash";
import {
  transactionDetailsData,
  routeProps,
  nodeRegions,
  error,
  timeScale,
  transaction,
  transactionFingerprintId,
} from "./transactionDetails.fixture";
import { SqlStatsSortOptions } from "src/api/statementsApi";

import { TransactionDetails } from ".";
import moment from "moment";

const lastUpdated = moment.utc();

storiesOf("Transactions Details", module)
  .addDecorator(storyFn => <MemoryRouter>{storyFn()}</MemoryRouter>)
  .addDecorator(storyFn => (
    <div style={{ backgroundColor: "#F5F7FA" }}>{storyFn()}</div>
  ))
  .add("with data", () => (
    <TransactionDetails
      {...routeProps}
      timeScale={timeScale}
      transactionFingerprintId={transactionFingerprintId.toString()}
      transaction={transaction}
      isLoading={false}
      statements={transactionDetailsData.statements}
      nodeRegions={nodeRegions}
      isTenant={false}
      hasViewActivityRedactedRole={false}
      refreshData={noop}
      refreshUserSQLRoles={noop}
      onTimeScaleChange={noop}
      refreshNodes={noop}
      lastUpdated={lastUpdated}
      isDataValid={true}
      limit={100}
      reqSortSetting={SqlStatsSortOptions.EXECUTION_COUNT}
    />
  ))
  .add("with loading indicator", () => (
    <TransactionDetails
      {...routeProps}
      timeScale={timeScale}
      transactionFingerprintId={transactionFingerprintId.toString()}
      transaction={null}
      isLoading={true}
      statements={undefined}
      nodeRegions={nodeRegions}
      isTenant={false}
      hasViewActivityRedactedRole={false}
      refreshData={noop}
      refreshUserSQLRoles={noop}
      onTimeScaleChange={noop}
      refreshNodes={noop}
      lastUpdated={lastUpdated}
      isDataValid={true}
      limit={100}
      reqSortSetting={SqlStatsSortOptions.EXECUTION_COUNT}
    />
  ))
  .add("with error alert", () => (
    <TransactionDetails
      {...routeProps}
      timeScale={timeScale}
      transactionFingerprintId={undefined}
      transaction={undefined}
      isLoading={false}
      statements={undefined}
      nodeRegions={nodeRegions}
      error={error}
      isTenant={false}
      hasViewActivityRedactedRole={false}
      refreshData={noop}
      refreshUserSQLRoles={noop}
      onTimeScaleChange={noop}
      refreshNodes={noop}
      lastUpdated={lastUpdated}
      isDataValid={false}
      limit={100}
      reqSortSetting={SqlStatsSortOptions.EXECUTION_COUNT}
    />
  ))
  .add("No data for this time frame; no cached transaction text", () => {
    return (
      <TransactionDetails
        {...routeProps}
        timeScale={timeScale}
        transactionFingerprintId={transactionFingerprintId.toString()}
        transaction={undefined}
        isLoading={false}
        statements={transactionDetailsData.statements}
        nodeRegions={nodeRegions}
        isTenant={false}
        hasViewActivityRedactedRole={false}
        refreshData={noop}
        refreshUserSQLRoles={noop}
        onTimeScaleChange={noop}
        refreshNodes={noop}
        lastUpdated={lastUpdated}
        isDataValid={true}
        limit={100}
        reqSortSetting={SqlStatsSortOptions.EXECUTION_COUNT}
      />
    );
  });
