// Copyright 2019 The Cockroach Authors.
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

import * as React from "react";
import { Helmet } from "react-helmet";
import { connect } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { Moment } from "moment";
import _ from "lodash";

import { AdminUIState } from "src/redux/state";
import { nodesSummarySelector } from "src/redux/nodes";
import { refreshLiveness, refreshNodes } from "src/redux/apiReducers";
import { cockroach } from "src/js/protos";
import MembershipStatus = cockroach.kv.kvserver.liveness.livenesspb.MembershipStatus;
import { LocalSetting } from "src/redux/localsettings";

import "./decommissionedNodeHistory.styl";
import { Text } from "src/components";
import {
  ColumnsConfig,
  Table,
  SortSetting,
  util,
} from "@cockroachlabs/cluster-ui";
import { createSelector } from "reselect";

const decommissionedNodesSortSetting = new LocalSetting<
  AdminUIState,
  SortSetting
>("nodes/decommissioned_sort_setting", s => s.localSettings);

interface DecommissionedNodeStatusRow {
  key: string;
  nodeId: number;
  decommissionedDate: Moment;
}

export interface DecommissionedNodeHistoryProps {
  refreshNodes: typeof refreshNodes;
  refreshLiveness: typeof refreshLiveness;
  dataSource: DecommissionedNodeStatusRow[];
}

const sortByNodeId = (
  a: DecommissionedNodeStatusRow,
  b: DecommissionedNodeStatusRow,
) => {
  if (a.nodeId < b.nodeId) {
    return -1;
  }
  if (a.nodeId > b.nodeId) {
    return 1;
  }
  return 0;
};

const sortByDecommissioningDate = (
  a: DecommissionedNodeStatusRow,
  b: DecommissionedNodeStatusRow,
) => {
  if (a.decommissionedDate.isBefore(b.decommissionedDate)) {
    return -1;
  }
  if (a.decommissionedDate.isAfter(b.decommissionedDate)) {
    return 1;
  }
  return 0;
};

export class DecommissionedNodeHistory extends React.Component<
  DecommissionedNodeHistoryProps
> {
  columns: ColumnsConfig<DecommissionedNodeStatusRow> = [
    {
      key: "id",
      title: "ID",
      sorter: sortByNodeId,
      render: (_text, record) => <Text>{`n${record.nodeId}`}</Text>,
    },
    {
      key: "decommissionedOn",
      title: "Decommissioned On",
      sorter: sortByDecommissioningDate,
      render: (_text, record) => {
        return record.decommissionedDate.format(util.DATE_FORMAT_24_UTC);
      },
    },
  ];

  componentDidMount() {
    this.props.refreshNodes();
    this.props.refreshLiveness();
  }

  componentDidUpdate() {
    this.props.refreshNodes();
    this.props.refreshLiveness();
  }

  render() {
    const { dataSource } = this.props;

    return (
      <section className="section">
        <Helmet title="Decommissioned Node History | Debug" />
        <h1 className="base-heading title">Decommissioned Node History</h1>
        <div>
          <Table
            dataSource={dataSource}
            columns={this.columns}
            noDataMessage="There are no decommissioned nodes in this cluster."
          />
        </div>
      </section>
    );
  }
}

const decommissionedNodesTableData = createSelector(
  nodesSummarySelector,
  (nodesSummary): DecommissionedNodeStatusRow[] => {
    const getDecommissionedTime = (nodeId: number) => {
      const liveness = nodesSummary.livenessByNodeID[nodeId];
      if (!liveness) {
        return undefined;
      }
      const deadTime = liveness.expiration.wall_time;
      return util.LongToMoment(deadTime);
    };

    const decommissionedNodes = Object.values(
      nodesSummary.livenessByNodeID,
    ).filter(liveness => {
      return liveness?.membership === MembershipStatus.DECOMMISSIONED;
    });

    const data = _.chain(decommissionedNodes)
      .orderBy([liveness => getDecommissionedTime(liveness.node_id)], ["desc"])
      .map((liveness, idx: number) => {
        const { node_id } = liveness;
        return {
          key: `${idx}`,
          nodeId: node_id,
          decommissionedDate: getDecommissionedTime(node_id),
        };
      })
      .value();

    return data;
  },
);

const mapStateToProps = (state: AdminUIState, _: RouteComponentProps) => ({
  dataSource: decommissionedNodesTableData(state),
});

const mapDispatchToProps = {
  refreshNodes,
  refreshLiveness,
  setSort: decommissionedNodesSortSetting.set,
};

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(DecommissionedNodeHistory),
);
