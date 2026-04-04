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

// createNodesByRegionMap creates a mapping of regions to nodes,
// based on the nodes list provided and nodeRegions, which is a full
// list of node id to the region it resides.
import * as protos from "src/js/protos";
import Long from "long";

type Timestamp = protos.google.protobuf.ITimestamp;

export function createNodesByRegionMap(
  nodes: number[],
  nodeRegions: Record<string, string>,
): Record<string, number[]> {
  const nodesByRegionMap: Record<string, number[]> = {};
  nodes.forEach((node: number) => {
    const region: string = nodeRegions[node.toString()];
    if (nodesByRegionMap[region] == null) {
      nodesByRegionMap[region] = [];
    }
    nodesByRegionMap[region].push(node);
  });
  return nodesByRegionMap;
}

// nodesByRegionMapToString converts a map of regions to node ids,
// ordered by region name, e.g. converts:
// { regionA: [1, 2], regionB: [2, 3] }
// to:
// regionA(n1, n2), regionB(n2,n3), ...
export function nodesByRegionMapToString(
  nodesByRegion: Record<string, number[]>,
): string {
  // Sort the nodes on each key.
  const regions = Object.keys(nodesByRegion).sort();
  regions.forEach((region: string) => {
    nodesByRegion[region].sort();
  });

  return regions
    .map((region: string) => {
      const nodes = nodesByRegion[region];
      return `${region}(${nodes.map(id => `n${id}`).join(",")})`;
    })
    .join(", ");
}

// getNodesByRegionString converts a list of node ids and map of
// node ids to region to a string of node ids by region, ordered
// by region name, e.g.
// regionA(n1, n2), regionB(n2,n3), ...
export function getNodesByRegionString(
  nodes: number[],
  nodeRegions: Record<string, string>,
): string {
  return nodesByRegionMapToString(createNodesByRegionMap(nodes, nodeRegions));
}

// makeTimestamp converts a string to a google.protobuf.Timestamp object.
export function makeTimestamp(date: string): Timestamp {
  return new protos.google.protobuf.Timestamp({
    seconds: new Long(Date.parse(date) * 1e-3),
  });
}
