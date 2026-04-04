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
import { Route, Switch, Redirect } from "react-router-dom";

import { NodesOverview } from "src/views/cluster/containers/nodesOverview";
import ClusterOverview from "src/views/cluster/containers/clusterOverview";

class NodesWrapper extends React.Component<{}, {}> {
  render() {
    return (
      <div
        style={{
          paddingTop: 12,
          width: "100%",
          height: "100%",
          overflow: "auto",
        }}
      >
        <NodesOverview />
      </div>
    );
  }
}

export default function createClusterOverviewRoutes(): JSX.Element {
  return (
    <Route path="/overview">
      <ClusterOverview>
        <Switch>
          <Redirect exact from="/overview" to="/overview/list" />
          <Route path="/overview/list" component={NodesWrapper} />
        </Switch>
      </ClusterOverview>
    </Route>
  );
}
