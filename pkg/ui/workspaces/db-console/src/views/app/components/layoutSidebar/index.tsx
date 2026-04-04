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
import { connect } from "react-redux";
import { Link, withRouter, RouteComponentProps } from "react-router-dom";

import { SideNavigation } from "src/components";
import "./navigation-bar.styl";
import { AdminUIState } from "src/redux/state";
import { isSingleNodeCluster } from "src/redux/nodes";

interface RouteParam {
  path: string;
  text: string;
  activeFor: string[];
  // ignoreFor allows exclude specific paths from been recognized as active
  // even if path is matched with activeFor paths.
  ignoreFor?: string[];
  isHidden?: () => boolean;
}

type SidebarProps = RouteComponentProps & MapToStateProps;

/**
 * Sidebar represents the static navigation sidebar available on all pages. It
 * displays a number of graphic icons representing available pages; the icon of
 * the page which is currently active will be highlighted.
 */
export class Sidebar extends React.Component<SidebarProps> {
  readonly routes: RouteParam[] = [
    { path: "/overview", text: "Overview", activeFor: ["/node"] },
    { path: "/metrics", text: "Metrics", activeFor: [] },
    { path: "/databases", text: "Databases", activeFor: ["/database"] },
    {
      path: "/sql-activity",
      text: "SQL Activity",
      activeFor: ["/sql-activity", "/session", "/transaction", "/statement"],
    },
    {
      path: "/reports/network",
      text: "Network Latency",
      activeFor: ["/reports/network"],
      // Do not show Network Latency for single node cluster.
      isHidden: () => this.props.isSingleNodeCluster,
    },
    {
      path: "/hotranges",
      text: "Hot Ranges",
      activeFor: ["/hotranges"],
    },
    { path: "/jobs", text: "Jobs", activeFor: [] },
    {
      path: "/debug",
      text: "Advanced Debug",
      activeFor: ["/reports", "/data-distribution", "/raft"],
      ignoreFor: ["/reports/network"],
    },
  ];

  isActiveNavigationItem = (path: string): boolean => {
    const { pathname } = this.props.location;
    const { activeFor, ignoreFor = [] } = this.routes.find(
      route => route.path === path,
    );
    return [...activeFor, path].some(p => {
      const isMatchedToIgnoredPaths = ignoreFor.some(ignorePath =>
        pathname.startsWith(ignorePath),
      );
      const isMatchedToActiveFor = pathname.startsWith(p);
      return isMatchedToActiveFor && !isMatchedToIgnoredPaths;
    });
  };

  render() {
    const navigationItems = this.routes
      .filter(route => !route.isHidden || !route.isHidden())
      .map(({ path, text }, idx) => (
        <SideNavigation.Item
          isActive={this.isActiveNavigationItem(path)}
          key={idx}
        >
          <Link to={path}>{text}</Link>
        </SideNavigation.Item>
      ));
    return <SideNavigation>{navigationItems}</SideNavigation>;
  }
}

interface MapToStateProps {
  isSingleNodeCluster: boolean;
}

const mapStateToProps = (state: AdminUIState) => ({
  isSingleNodeCluster: isSingleNodeCluster(state),
});

export default withRouter(connect(mapStateToProps, null)(Sidebar));
