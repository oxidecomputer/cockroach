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

import { LineGraph } from "src/views/cluster/components/linegraph";
import { Metric, Axis } from "src/views/shared/components/metricQuery";

import {
  GraphDashboardProps,
  nodeDisplayName,
  storeIDsForNode,
} from "./dashboardUtils";
import { AvailableDiscCapacityGraphTooltip } from "src/views/cluster/containers/nodeGraphs/dashboards/graphTooltips";
import { AxisUnits } from "@cockroachlabs/cluster-ui";

// TODO(vilterp): tooltips

export default function(props: GraphDashboardProps) {
  const {
    nodeIDs,
    nodeDisplayNameByID,
    storeIDsByNodeID,
    nodeSources,
    storeSources,
    tooltipSelection,
  } = props;

  return [
    <LineGraph title="CPU Percent" sources={nodeSources}>
      <Axis units={AxisUnits.Percentage} label="CPU">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.cpu.combined.percent-normalized"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph
      title="Memory Usage"
      sources={nodeSources}
      tooltip={<div>Memory in use {tooltipSelection}</div>}
    >
      <Axis units={AxisUnits.Bytes} label="memory usage">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.rss"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph title="Disk Read MiB/s" sources={nodeSources}>
      <Axis units={AxisUnits.Bytes} label="bytes">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.host.disk.read.bytes"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
            nonNegativeRate
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph title="Disk Write MiB/s" sources={nodeSources}>
      <Axis units={AxisUnits.Bytes} label="bytes">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.host.disk.write.bytes"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
            nonNegativeRate
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph title="Disk Read IOPS" sources={nodeSources}>
      <Axis units={AxisUnits.Count} label="IOPS">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.host.disk.read.count"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
            nonNegativeRate
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph title="Disk Write IOPS" sources={nodeSources}>
      <Axis units={AxisUnits.Count} label="IOPS">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.host.disk.write.count"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
            nonNegativeRate
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph title="Disk Ops In Progress" sources={nodeSources}>
      <Axis units={AxisUnits.Count} label="Ops">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.host.disk.iopsinprogress"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph
      title="Available Disk Capacity"
      sources={storeSources}
      tooltip={<AvailableDiscCapacityGraphTooltip />}
    >
      <Axis units={AxisUnits.Bytes} label="capacity">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.store.capacity.available"
            sources={storeIDsForNode(storeIDsByNodeID, nid)}
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph title="Network Bytes Received" sources={nodeSources}>
      <Axis units={AxisUnits.Bytes} label="bytes">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.host.net.recv.bytes"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
            nonNegativeRate
          />
        ))}
      </Axis>
    </LineGraph>,

    <LineGraph title="Network Bytes Sent" sources={nodeSources}>
      <Axis units={AxisUnits.Bytes} label="bytes">
        {nodeIDs.map(nid => (
          <Metric
            name="cr.node.sys.host.net.send.bytes"
            title={nodeDisplayName(nodeDisplayNameByID, nid)}
            sources={[nid]}
            nonNegativeRate
          />
        ))}
      </Axis>
    </LineGraph>,
  ];
}
