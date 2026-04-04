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

import { createSelector } from "reselect";
import { IndexDetailsPageData, util } from "@cockroachlabs/cluster-ui";
import { AdminUIState } from "src/redux/state";
import { RouteComponentProps } from "react-router";
import { getMatchParamByName } from "src/util/query";
import {
  databaseNameAttr,
  tableNameAttr,
  indexNameAttr,
} from "src/util/constants";
import {
  generateTableID,
  refreshIndexStats,
  refreshNodes,
  refreshUserSQLRoles,
} from "src/redux/apiReducers";
import { resetIndexUsageStatsAction } from "src/redux/indexUsageStats";
import { longToInt } from "src/util/fixLong";
import { cockroach } from "src/js/protos";
import TableIndexStatsRequest = cockroach.server.serverpb.TableIndexStatsRequest;
import { selectHasAdminRole } from "src/redux/user";

export const mapStateToProps = createSelector(
  (_state: AdminUIState, props: RouteComponentProps): string =>
    getMatchParamByName(props.match, databaseNameAttr),
  (_state: AdminUIState, props: RouteComponentProps): string =>
    getMatchParamByName(props.match, tableNameAttr),
  (_state: AdminUIState, props: RouteComponentProps): string =>
    getMatchParamByName(props.match, indexNameAttr),
  state => state.cachedData.indexStats,
  state => selectHasAdminRole(state),
  (database, table, index, indexStats, hasAdminRole): IndexDetailsPageData => {
    const stats = indexStats[generateTableID(database, table)];
    const details = stats?.data?.statistics.filter(
      stat => stat.index_name === index, // index names must be unique for a table
    )[0];
    return {
      databaseName: database,
      tableName: table,
      indexName: index,
      hasAdminRole: hasAdminRole,
      details: {
        loading: !!stats?.inFlight,
        loaded: !!stats?.valid,
        createStatement: details?.create_statement || "",
        totalReads:
          longToInt(details?.statistics?.stats?.total_read_count) || 0,
        lastRead: util.TimestampToMoment(details?.statistics?.stats?.last_read),
        lastReset: util.TimestampToMoment(stats?.data?.last_reset),
      },
    };
  },
);

export const mapDispatchToProps = {
  refreshIndexStats: (database: string, table: string) => {
    return refreshIndexStats(new TableIndexStatsRequest({ database, table }));
  },
  resetIndexUsageStats: resetIndexUsageStatsAction,
  refreshNodes,
  refreshUserSQLRoles,
};
