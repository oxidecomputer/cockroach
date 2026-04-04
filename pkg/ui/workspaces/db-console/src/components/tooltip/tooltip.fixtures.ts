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

import { NodeStatusRow } from "src/views/cluster/containers/nodesOverview";

export const nodeLocalityFixture: NodeStatusRow = {
  key: "-0",
  nodeId: 1,
  nodeName: "localhost:26257",
  uptime: "3 hours",
  replicas: 34,
  usedCapacity: 135351337,
  availableCapacity: 108590390313,
  usedMemory: 151085056,
  availableMemory: 8589934592,
  numCpus: 4,
  version: "v20.2.0-alpha.1-1355-ga0123f1bc0",
  status: 3,
  tiers: [
    {
      key: "region",
      value: "gcp-us-east1",
    },
  ],
  region: "gcp-us-east1",
};
