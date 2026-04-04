// Copyright 2018 The Cockroach Authors.
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
import { Helmet } from "react-helmet";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { connect } from "react-redux";

import NavigationBar from "src/views/app/components/layoutSidebar";
import ErrorBoundary from "src/views/app/components/errorMessage/errorBoundary";
import TimeWindowManager from "src/views/app/containers/metricsTimeManager";
import AlertBanner from "src/views/app/containers/alertBanner";
import RequireLogin from "src/views/login/requireLogin";
import {
  clusterIdSelector,
  clusterNameSelector,
  singleVersionSelector,
} from "src/redux/nodes";
import { AdminUIState } from "src/redux/state";
import LoginIndicator from "src/views/app/components/loginIndicator";
import {
  GlobalNavigation,
  CockroachLabsLockupIcon,
  Left,
  Right,
  PageHeader,
  Text,
  TextTypes,
  Badge,
} from "src/components";

import "./layout.styl";
import "./layoutPanel.styl";

export interface LayoutProps {
  clusterName: string;
  clusterVersion: string;
  clusterId: string;
}

/**
 * Defines the main layout of all admin ui pages. This includes static
 * navigation bars and footers which should be present on every page.
 *
 * Individual pages provide their content via react-router.
 */
class Layout extends React.Component<LayoutProps & RouteComponentProps> {
  contentRef = React.createRef<HTMLDivElement>();

  componentDidUpdate(prevProps: RouteComponentProps) {
    // `react-router` doesn't handle scroll restoration (https://reactrouter.com/react-router/web/guides/scroll-restoration)
    // and when location changed with react-router's Link it preserves scroll position whenever it is.
    // AdminUI layout keeps left and top panels have fixed position on a screen and has internal scrolling for content div
    // element which has to be scrolled back on top with navigation change.
    if (this.props.location.pathname !== prevProps.location.pathname) {
      this.contentRef.current.scrollTo(0, 0);
    }
  }

  render() {
    const { clusterName, clusterVersion, clusterId } = this.props;
    return (
      <RequireLogin>
        <Helmet
          titleTemplate="%s | Cockroach Console"
          defaultTitle="Cockroach Console"
        />
        <TimeWindowManager />
        <AlertBanner />
        <div className="layout-panel">
          <div className="layout-panel__header">
            <GlobalNavigation>
              <Left>
                <CockroachLabsLockupIcon height={26} />
              </Left>
              <Right>
                <LoginIndicator />
              </Right>
            </GlobalNavigation>
          </div>
          <div className="layout-panel__navigation-bar">
            <PageHeader>
              <Text textType={TextTypes.Heading2} noWrap>
                {clusterName || `Cluster id: ${clusterId || ""}`}
              </Text>
              <Badge text={clusterVersion} />
            </PageHeader>
          </div>
          <div className="layout-panel__body">
            <div className="layout-panel__sidebar">
              <NavigationBar />
            </div>
            <div ref={this.contentRef} className="layout-panel__content">
              <ErrorBoundary key={this.props.location.pathname}>
                {this.props.children}
              </ErrorBoundary>
            </div>
          </div>
        </div>
      </RequireLogin>
    );
  }
}

const mapStateToProps = (state: AdminUIState) => {
  return {
    clusterName: clusterNameSelector(state),
    clusterVersion: singleVersionSelector(state),
    clusterId: clusterIdSelector(state),
  };
};

export default withRouter(connect(mapStateToProps)(Layout));
