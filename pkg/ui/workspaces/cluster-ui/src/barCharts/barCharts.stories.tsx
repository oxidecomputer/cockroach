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
import { storiesOf, DecoratorFn } from "@storybook/react";

import {
  countBarChart,
  bytesReadBarChart,
  latencyBarChart,
  maxMemUsageBarChart,
  networkBytesBarChart,
  retryBarChart,
} from "./barCharts";
import statementsPagePropsFixture from "src/statementsPage/statementsPage.fixture";
import Long from "long";

const { statements } = statementsPagePropsFixture;

const withinColumn = (width = "150px"): DecoratorFn => storyFn => {
  const rowStyle = {
    borderTop: "1px solid #e7ecf3",
    borderBottom: "1px solid #e7ecf3",
  };

  const cellStyle = {
    width: "190px",
    padding: "10px 20px",
  };

  return (
    <table>
      <tbody>
        <tr style={rowStyle}>
          <td style={cellStyle}>
            <div style={{ width }}>{storyFn()}</div>
          </td>
        </tr>
      </tbody>
    </table>
  );
};

storiesOf("BarCharts", module)
  .add("countBarChart", () => {
    const chartFactory = countBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("bytesReadBarChart", () => {
    const chartFactory = bytesReadBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("latencyBarChart", () => {
    const chartFactory = latencyBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("maxMemUsageBarChart", () => {
    const chartFactory = maxMemUsageBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("networkBytesBarChart", () => {
    const chartFactory = networkBytesBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("retryBarChart", () => {
    const chartFactory = retryBarChart(statements);
    return chartFactory(statements[0]);
  });

storiesOf("BarCharts/within column (150px)", module)
  .addDecorator(withinColumn())
  .add("countBarChart", () => {
    const chartFactory = countBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("bytesReadBarChart", () => {
    const chartFactory = bytesReadBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("latencyBarChart", () => {
    const chartFactory = latencyBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("maxMemUsageBarChart", () => {
    const chartFactory = maxMemUsageBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("retryBarChart", () => {
    const chartFactory = retryBarChart(statements);
    return chartFactory(statements[0]);
  })
  .add("empty retryBarChart", () => {
    const withoutRetries = statements.map(s => ({
      ...s,
      stats: {
        ...s.stats,
        count: Long.fromNumber(0),
        first_attempt_count: Long.fromNumber(0),
        max_retries: Long.fromNumber(0),
      },
    }));
    const chartFactory = retryBarChart(withoutRetries);
    return chartFactory(withoutRetries[0]);
  });
