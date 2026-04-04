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

import React, { useEffect, useRef } from "react";
import classNames from "classnames/bind";
import { getStackedBarOpts, stack } from "./bars";
import uPlot, { AlignedData } from "uplot";
import styles from "./bargraph.module.scss";
import { Visualization } from "../visualization";
import {
  AxisUnits,
  calculateXAxisDomainBarChart,
  calculateYAxisDomain,
} from "../utils/domain";
import { Options } from "uplot";

const cx = classNames.bind(styles);

export type BarGraphTimeSeriesProps = {
  alignedData?: AlignedData;
  colourPalette?: string[]; // Series colour palette.
  preCalcGraphSize?: boolean;
  title: string;
  tooltip?: React.ReactNode;
  uPlotOptions: Partial<Options>;
  yAxisUnits: AxisUnits;
};

// Currently this component only supports stacked multi-series bars.
export const BarGraphTimeSeries: React.FC<BarGraphTimeSeriesProps> = ({
  alignedData,
  colourPalette,
  preCalcGraphSize = true,
  title,
  tooltip,
  uPlotOptions,
  yAxisUnits,
}) => {
  const graphRef = useRef<HTMLDivElement>(null);
  const samplingIntervalMillis =
    alignedData[0].length > 1 ? alignedData[0][1] - alignedData[0][0] : 1e3;

  useEffect(() => {
    if (!alignedData) return;

    const xAxisDomain = calculateXAxisDomainBarChart(
      alignedData[0][0], // startMillis
      alignedData[0][alignedData[0].length - 1], // endMillis
      samplingIntervalMillis,
    );

    const stackedData = stack(alignedData, () => false);

    const allYDomainPoints: number[] = [];
    stackedData.slice(1).forEach(points => allYDomainPoints.push(...points));
    const yAxisDomain = calculateYAxisDomain(yAxisUnits, allYDomainPoints);

    const opts = getStackedBarOpts(
      alignedData,
      uPlotOptions,
      xAxisDomain,
      yAxisDomain,
      yAxisUnits,
      colourPalette,
    );

    const plot = new uPlot(opts, stackedData, graphRef.current);

    return () => {
      plot?.destroy();
    };
  }, [
    alignedData,
    colourPalette,
    uPlotOptions,
    yAxisUnits,
    samplingIntervalMillis,
  ]);

  return (
    <Visualization
      title={title}
      loading={!alignedData}
      preCalcGraphSize={preCalcGraphSize}
      tooltip={tooltip}
    >
      <div className={cx("bargraph")}>
        <div ref={graphRef} />
      </div>
    </Visualization>
  );
};
