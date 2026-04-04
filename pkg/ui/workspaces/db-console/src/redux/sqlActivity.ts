// Copyright 2022 The Cockroach Authors.
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
import { Action } from "redux";
import _ from "lodash";
import { PayloadAction } from "src/interfaces/action";

/**
 * SqlActivityState maintains a MetricQuerySet collection, along with some
 * metadata relevant to server queries.
 */
export class SqlActivityState {
  /**
   * Caches the latest query text associated with the statement fingerprint in the URL, to preserve this data even if
   * the time frame changes such that there is no longer data for this statement fingerprint in the selected time frame.
   */
  statementDetailsLatestQuery: string;
  /**
   * Caches the latest formatted query associated with the statement fingerprint in the URL, to preserve this data even if
   * the time frame changes such that there is no longer data for this statement fingerprint in the selected time frame.
   */
  statementDetailsLatestFormattedQuery: string;
}

const SET_STATEMENT_DETAILS_LATEST_QUERY =
  "cockroachui/sqlActivity/SET_STATEMENT_DETAILS_LATEST_QUERY";
const SET_STATEMENT_DETAILS_LATEST_FORMATTED_QUERY =
  "cockroachui/sqlActivity/SET_STATEMENT_DETAILS_LATEST_FORMATTED_QUERY";

export function statementDetailsLatestQueryAction(
  query: string,
): PayloadAction<string> {
  return {
    type: SET_STATEMENT_DETAILS_LATEST_QUERY,
    payload: query,
  };
}
export function statementDetailsLatestFormattedQueryAction(
  formattedQuery: string,
): PayloadAction<string> {
  return {
    type: SET_STATEMENT_DETAILS_LATEST_FORMATTED_QUERY,
    payload: formattedQuery,
  };
}

/**
 * The metrics reducer accepts events for individual MetricQuery objects,
 * dispatching them based on ID. It also accepts actions which indicate the
 * state of the connection to the server.
 */
export function sqlActivityReducer(
  state: SqlActivityState = new SqlActivityState(),
  action: Action,
): SqlActivityState {
  switch (action.type) {
    case SET_STATEMENT_DETAILS_LATEST_QUERY: {
      const { payload: query } = action as PayloadAction<string>;
      state = _.clone(state);
      state.statementDetailsLatestQuery = query;
      return state;
    }
    case SET_STATEMENT_DETAILS_LATEST_FORMATTED_QUERY: {
      const { payload: formattedQuery } = action as PayloadAction<string>;
      state = _.clone(state);
      state.statementDetailsLatestFormattedQuery = formattedQuery;
      return state;
    }
    default:
      return state;
  }
}
