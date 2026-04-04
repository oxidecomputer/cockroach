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

import React, { useCallback, useEffect, useState } from "react";
import { RouteComponentProps, withRouter } from "react-router-dom";
import moment from "moment";
import { Button } from "@cockroachlabs/ui-components";
import { cockroach } from "src/js/protos";
import { getHotRanges } from "src/util/api";

type HotRangesProps = RouteComponentProps<{ node_id: string }>;

const HotRanges = (props: HotRangesProps) => {
  const nodeIdParam = props.match.params["node_id"];
  const [nodeId, setNodeId] = useState(nodeIdParam);
  const [time, setTime] = useState<moment.Moment>(moment());
  const [hotRanges, setHotRanges] = useState<
    cockroach.server.serverpb.HotRangesResponseV2["ranges"]
  >([]);
  const [pageToken, setPageToken] = useState<string>("");
  const pageSize = 50;

  const refreshHotRanges = useCallback(() => {
    setHotRanges([]);
    setPageToken("");
  }, []);

  useEffect(() => {
    const request = new cockroach.server.serverpb.HotRangesRequest({
      node_id: nodeId,
      page_size: pageSize,
      page_token: pageToken,
    });
    getHotRanges(request).then(response => {
      if (response.ranges.length == 0) {
        return;
      }
      setPageToken(response.next_page_token);
      setHotRanges([...hotRanges, ...response.ranges]);
      setTime(moment());
    });
    // Avoid dispatching request when `hotRanges` list is updated.
    // This effect should be triggered only when pageToken is changed.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageToken]);

  useEffect(() => {
    setNodeId(nodeIdParam);
  }, [nodeIdParam]);
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
      }}
    >
      <span>{`Node ID: ${nodeId ?? "All nodes"}`}</span>
      <span>{`Time: ${time.toISOString()}`}</span>
      <Button onClick={refreshHotRanges} intent={"secondary"}>
        Refresh
      </Button>
      <pre className="state-json-box">{JSON.stringify(hotRanges, null, 2)}</pre>
    </div>
  );
};

export default withRouter(HotRanges);
