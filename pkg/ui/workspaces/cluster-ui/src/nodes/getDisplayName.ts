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

import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { NoConnection } from "./noConnection";

type INodeStatus = cockroach.server.status.statuspb.INodeStatus;

const LivenessStatus =
  cockroach.kv.kvserver.liveness.livenesspb.NodeLivenessStatus;

function isNoConnection(
  node: INodeStatus | NoConnection,
): node is NoConnection {
  return (
    (node as NoConnection).to !== undefined &&
    (node as NoConnection).from !== undefined
  );
}

export function getDisplayName(
  node: INodeStatus | NoConnection,
  livenessStatus = LivenessStatus.NODE_STATUS_LIVE,
) {
  const decommissionedString =
    livenessStatus === LivenessStatus.NODE_STATUS_DECOMMISSIONED
      ? "[decommissioned] "
      : "";

  if (isNoConnection(node)) {
    return `${decommissionedString}(n${node.from.nodeID})`;
  }
  // as the only other type possible right now is INodeStatus we don't have a type guard for that
  return `${decommissionedString}(n${node.desc.node_id}) ${node.desc.address.address_field}`;
}
