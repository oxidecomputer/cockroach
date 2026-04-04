// Copyright 2020 The Cockroach Authors.
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

import * as ClusterTooltips from "src/views/cluster/containers/clusterOverview/tooltips";
import * as NodeOverviewTooltips from "src/views/cluster/containers/nodeOverview/tooltips";
import * as GraphTooltips from "src/views/cluster/containers/nodeGraphs/dashboards/graphTooltips";
import { ToolTipWrapper } from "src/views/shared/components/toolTip";
import {
  plainNodeTooltips,
  getNodeStatusDescription,
  getStatusDescription,
  NodeLocalityColumn,
} from "src/views/cluster/containers/nodesOverview/tooltips";
import { AggregatedNodeStatus } from "src/views/cluster/containers/nodesOverview";
import { LivenessStatus } from "src/redux/nodes";
import { Tooltip } from "src/components/tooltip/tooltip";
import { styledWrapper } from "src/util/decorators";
import { nodeLocalityFixture } from "src/components/tooltip/tooltip.fixtures";

const triggerStyle: React.CSSProperties = {
  width: "300px",
  marginBottom: "300px",
};

const graphTooltipsStyle: React.CSSProperties = {
  width: "450px",
  textAlign: "center",
};

const wrapperStyle: React.CSSProperties = {
  padding: "24px",
  display: "flex",
  flexWrap: "wrap",
};

const TooltipTrigger = (props: {
  name?: string;
  children?: React.ReactNode;
}) => (
  <button style={triggerStyle}>
    {props.name || null}
    {props.children || null}
  </button>
);

const tooltipsStack = (components: object) => (
  <>
    {Object.values(components).map((Item, idx) => (
      <Item visible={true} key={idx}>
        {TooltipTrigger({ name: Item.name })}
      </Item>
    ))}
  </>
);

const statusTooltipsStack = (statusNames: object, descriptionGetter: any) => (
  <>
    {Object.keys(statusNames)
      .filter(status => isNaN(Number(status)))
      .map((status: keyof typeof statusNames, idx) => (
        <Tooltip
          key={idx}
          title={descriptionGetter(statusNames[status])}
          visible={true}
          placement="bottom"
        >
          {TooltipTrigger({ name: status })}
        </Tooltip>
      ))}
  </>
);

const graphTooltipsStack = (components: object) => (
  <>
    {Object.values(components).map((Item, idx) => (
      <div style={graphTooltipsStyle}>
        <ToolTipWrapper
          key={idx}
          text={<Item tooltipSelection={"on <node>"} />}
          visible={true}
        >
          {TooltipTrigger({ name: Item.name })}
        </ToolTipWrapper>
      </div>
    ))}
  </>
);

storiesOf("Tooltips/Cluster Overview Page", module)
  .addDecorator(styledWrapper(wrapperStyle))
  .add("Summary tooltips", () => tooltipsStack(ClusterTooltips))
  .add("Node overview tooltips", () => tooltipsStack(NodeOverviewTooltips))
  .add("Node List tooltips", () => (
    <>
      {tooltipsStack(plainNodeTooltips)}
      {TooltipTrigger({
        children: (
          <NodeLocalityColumn record={nodeLocalityFixture} visible={true} />
        ),
      })}
    </>
  ))
  .add("Node Status tooltips", () => (
    <>
      {statusTooltipsStack(AggregatedNodeStatus, getNodeStatusDescription)}
      {statusTooltipsStack(LivenessStatus, getStatusDescription)}
    </>
  ));
storiesOf("Tooltips/Metrics Page", module)
  .addDecorator(styledWrapper(wrapperStyle))
  .add("Graph tooltips", () => graphTooltipsStack(GraphTooltips));
