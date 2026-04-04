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

import React from "react";
import { cockroach } from "@cockroachlabs/crdb-protobuf-client";
import { fetchData } from "src/api";

type StatementsResponse = cockroach.server.serverpb.StatementsResponse;

export const DemoFetch: React.FC = () => {
  const [response, setResponse] = React.useState<StatementsResponse>(null);

  React.useEffect(() => {
    fetchData(
      cockroach.server.serverpb.StatementsResponse,
      "_status/statements",
    ).then(setResponse);
  }, []);

  return <code>{JSON.stringify(response).slice(0, 255)}...</code>;
};
