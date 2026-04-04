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
import { Tooltip, Anchor } from "src/components";
import {
  keyValuePairs,
  writeIntents,
  metaRanges,
  clusterStore,
  capacityMetrics,
} from "src/util/docs";
import { TooltipProps } from "src/components/tooltip/tooltip";

export const LiveBytesTooltip: React.FC<TooltipProps & {
  nodeName: string;
}> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>
          Number of logical bytes stored in live &nbsp;
          <Anchor href={keyValuePairs} target="_blank">
            key-value pairs
          </Anchor>{" "}
          &nbsp; on node {props.nodeName || "NodeName"}.
        </p>
        <p>Live data excludes historical and deleted data.</p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);

export const KeyBytesTooltip: React.FC<TooltipProps & {
  nodeName: string;
}> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>
          Number of bytes stored in keys on node {props.nodeName || "NodeName"}.
        </p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);

export const ValueBytesTooltip: React.FC<TooltipProps & {
  nodeName: string;
}> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>
          Number of bytes stored in values on node{" "}
          {props.nodeName || "NodeName"}.
        </p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);

export const IntentBytesTooltip: React.FC<TooltipProps> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>
          Number of bytes stored in &nbsp;
          <Anchor href={writeIntents} target="_blank">
            write intents
          </Anchor>{" "}
          &nbsp; of uncommitted values.
        </p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);

export const SystemBytesTooltip: React.FC<TooltipProps & {
  nodeName: string;
}> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>
          Number of physical bytes stored in &nbsp;
          <Anchor href={metaRanges} target="_blank">
            system key-value pairs
          </Anchor>{" "}
          &nbsp; on node {props.nodeName || "NodeName"}.
        </p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);

export const NodeUsedCapacityTooltip: React.FC<TooltipProps & {
  nodeName: string;
}> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>
          Disk space in use by CockroachDB data on node{" "}
          {props.nodeName || "NodeName"}.
        </p>
        <p>
          This excludes the Cockroach binary, operating system, and other system
          files.
        </p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);

export const NodeAvailableCapacityTooltip: React.FC<TooltipProps & {
  nodeName: string;
}> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>
          Free disk space available to CockroachDB data on node{" "}
          {props.nodeName || "NodeName"}.
        </p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);

export const NodeMaximumCapacityTooltip: React.FC<TooltipProps & {
  nodeName: string;
}> = props => (
  <Tooltip
    {...props}
    placement="bottom"
    title={
      <div className="tooltip__table--title">
        <p>Maximum store size of node {props.nodeName || "NodeName"}.</p>
        <p>
          This value may be explicitly set per node using &nbsp;
          <Anchor href={clusterStore} target="_blank">
            --store
          </Anchor>{" "}
          &nbsp; If a store size has not been set, this metric displays the
          actual disk capacity.
        </p>
        <p>
          <Anchor href={capacityMetrics} target="_blank">
            How is this metric calculated?
          </Anchor>
        </p>
      </div>
    }
  >
    {props.children}
  </Tooltip>
);
