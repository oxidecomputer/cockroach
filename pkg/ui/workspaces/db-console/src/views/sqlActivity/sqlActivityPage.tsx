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

// All changes made on this file, should also be done on the equivalent
// file on managed-service repo.

import React, { useState, useEffect } from "react";
import Helmet from "react-helmet";
import { Tabs } from "antd";
import { commonStyles, util } from "@cockroachlabs/cluster-ui";
import SessionsPageConnected from "src/views/sessions/sessionsPage";
import TransactionsPageConnected from "src/views/transactions/transactionsPage";
import StatementsPageConnected from "src/views/statements/statementsPage";
import { RouteComponentProps } from "react-router-dom";

const { TabPane } = Tabs;

const SQLActivityPage = (props: RouteComponentProps) => {
  const defaultTab = util.queryByName(props.location, "tab") || "Statements";
  const [currentTab, setCurrentTab] = useState(defaultTab);

  const onTabChange = (tabId: string): void => {
    setCurrentTab(tabId);
    props.history.location.search = "";
    util.syncHistory({ tab: tabId }, props.history, true);
  };

  useEffect(() => {
    const queryTab = util.queryByName(props.location, "tab") || "Statements";
    if (queryTab !== currentTab) {
      setCurrentTab(queryTab);
    }
  }, [props.location, currentTab]);

  return (
    <div>
      <Helmet title={defaultTab} />
      <h3 className={commonStyles("base-heading")}>SQL Activity</h3>
      <Tabs
        defaultActiveKey={defaultTab}
        className={commonStyles("cockroach--tabs")}
        onChange={onTabChange}
        activeKey={currentTab}
      >
        <TabPane tab="Statements" key="Statements">
          <StatementsPageConnected />
        </TabPane>
        <TabPane tab="Transactions" key="Transactions">
          <TransactionsPageConnected />
        </TabPane>
        <TabPane tab="Sessions" key="Sessions">
          <SessionsPageConnected />
        </TabPane>
      </Tabs>
    </div>
  );
};

export default SQLActivityPage;
