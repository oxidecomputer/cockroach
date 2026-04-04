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
import { cloneDeep, noop, extend } from "lodash";
import {
  data,
  nodeRegions,
  columns,
  routeProps,
  timeScale,
  sortSetting,
  filters,
} from "./transactions.fixture";

import { TransactionsPage } from ".";
import { RequestError } from "../util";
import moment from "moment";
import { SqlStatsSortOptions } from "../api";

const getEmptyData = () =>
  extend({}, data, { transactions: [], statements: [] });

const lastUpdated = moment.utc();

const defaultReqParaProps = {
  limit: 100,
  reqSortSetting: SqlStatsSortOptions.PCT_RUNTIME,
  onChangeLimit: noop,
  onChangeReqSort: noop,
};

storiesOf("Transactions Page", module)
  .addDecorator(storyFn => <MemoryRouter>{storyFn()}</MemoryRouter>)
  .addDecorator(storyFn => (
    <div style={{ backgroundColor: "#F5F7FA" }}>{storyFn()}</div>
  ))
  .add("with data", () => (
    <TransactionsPage
      {...routeProps}
      columns={columns}
      data={data}
      timeScale={timeScale}
      filters={filters}
      nodeRegions={nodeRegions}
      hasAdminRole={true}
      onFilterChange={noop}
      onSortingChange={noop}
      refreshData={noop}
      refreshNodes={noop}
      refreshUserSQLRoles={noop}
      resetSQLStats={noop}
      search={""}
      sortSetting={sortSetting}
      lastUpdated={lastUpdated}
      isDataValid={true}
      isReqInFlight={false}
      {...defaultReqParaProps}
    />
  ))
  .add("without data", () => {
    return (
      <TransactionsPage
        {...routeProps}
        columns={columns}
        data={getEmptyData()}
        timeScale={timeScale}
        filters={filters}
        nodeRegions={nodeRegions}
        hasAdminRole={true}
        onFilterChange={noop}
        onSortingChange={noop}
        refreshData={noop}
        refreshNodes={noop}
        refreshUserSQLRoles={noop}
        resetSQLStats={noop}
        search={""}
        sortSetting={sortSetting}
        lastUpdated={lastUpdated}
        isDataValid={true}
        isReqInFlight={false}
        {...defaultReqParaProps}
      />
    );
  })
  .add("with empty search result", () => {
    const route = cloneDeep(routeProps);
    const { history } = route;
    const searchParams = new URLSearchParams(history.location.search);
    searchParams.set("q", "aaaaaaa");
    history.location.search = searchParams.toString();

    return (
      <TransactionsPage
        {...routeProps}
        columns={columns}
        data={getEmptyData()}
        timeScale={timeScale}
        filters={filters}
        history={history}
        nodeRegions={nodeRegions}
        hasAdminRole={true}
        onFilterChange={noop}
        onSortingChange={noop}
        refreshData={noop}
        refreshNodes={noop}
        refreshUserSQLRoles={noop}
        resetSQLStats={noop}
        search={""}
        sortSetting={sortSetting}
        lastUpdated={lastUpdated}
        isDataValid={true}
        isReqInFlight={false}
        {...defaultReqParaProps}
      />
    );
  })
  .add("with loading indicator", () => {
    return (
      <TransactionsPage
        {...routeProps}
        columns={columns}
        data={undefined}
        timeScale={timeScale}
        filters={filters}
        nodeRegions={nodeRegions}
        hasAdminRole={true}
        onFilterChange={noop}
        onSortingChange={noop}
        refreshData={noop}
        refreshNodes={noop}
        refreshUserSQLRoles={noop}
        resetSQLStats={noop}
        search={""}
        sortSetting={sortSetting}
        lastUpdated={lastUpdated}
        isDataValid={true}
        isReqInFlight={true}
        {...defaultReqParaProps}
      />
    );
  })
  .add("with error alert", () => {
    return (
      <TransactionsPage
        {...routeProps}
        columns={columns}
        data={undefined}
        timeScale={timeScale}
        error={
          new RequestError(
            "Forbidden",
            403,
            "this operation requires admin privilege",
          )
        }
        filters={filters}
        nodeRegions={nodeRegions}
        hasAdminRole={true}
        onFilterChange={noop}
        onSortingChange={noop}
        refreshData={noop}
        refreshNodes={noop}
        refreshUserSQLRoles={noop}
        resetSQLStats={noop}
        search={""}
        sortSetting={sortSetting}
        lastUpdated={lastUpdated}
        isDataValid={false}
        isReqInFlight={false}
        {...defaultReqParaProps}
      />
    );
  });
