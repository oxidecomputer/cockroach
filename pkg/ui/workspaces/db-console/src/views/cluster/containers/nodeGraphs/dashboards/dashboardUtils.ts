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

/**
 * GraphDashboardProps are the properties accepted by the renderable component
 * of each graph dashboard.
 */
export interface GraphDashboardProps {
  /**
   * List of node IDs which should be used in graphs which display a series per
   * node.
   */
  nodeIDs: string[];
  /**
   * List of nodes which should be queried for data. This will be empty if all
   * nodes should be queried.
   */
  nodeSources: string[];
  /**
   * List of stores which should be displayed in the dashboard. This will be
   * empty if all stores should be queried.
   */
  storeSources: string[];
  /**
   * tooltipSelection is a string used in tooltips to reference the currently
   * selected nodes. This is a prepositional phrase, currently either "across
   * all nodes" or "on node X".
   */
  tooltipSelection: string;

  nodeDisplayNameByID: {
    [key: string]: string;
  };

  storeIDsByNodeID: {
    [key: string]: string[];
  };
}

export function nodeDisplayName(
  nodeDisplayNameByID: { [nodeId: string]: string },
  nid: string,
): string {
  return nodeDisplayNameByID[nid] || "unknown node";
}

export function storeIDsForNode(
  storeIDsByNodeID: { [key: string]: string[] },
  nid: string,
): string[] {
  return storeIDsByNodeID[nid] || [];
}
