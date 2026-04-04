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

import React from "react";
import classNames from "classnames/bind";
import spinner from "src/assets/spinner.gif";
import { Tooltip } from "antd";

import styles from "./visualizations.module.scss";
const cx = classNames.bind(styles);

interface VisualizationProps {
  title: string;
  subtitle?: string;
  tooltip?: React.ReactNode;
  // If stale is true, the visualization is faded
  // and the icon is changed to a warning icon.
  stale?: boolean;
  // If loading is true a spinner is shown instead of the graph.
  loading?: boolean;
  preCalcGraphSize?: boolean;
  children: React.ReactNode;
}

/**
 * Visualization is a container for a variety of visual elements (such as
 * charts). It surrounds a visual element with some standard information, such
 * as a title and a tooltip icon.
 */
export const Visualization: React.FC<VisualizationProps> = ({
  title,
  subtitle,
  tooltip,
  stale,
  loading,
  preCalcGraphSize,
  children,
}) => {
  const chartTitle: React.ReactNode = (
    <div>
      <span
        className={cx(
          "visualization-title",
          tooltip && "visualization-underline",
        )}
      >
        {title}
      </span>
      {subtitle && (
        <span className={cx("visualization-subtitle")}>{subtitle}</span>
      )}
    </div>
  );

  const tooltipNode: React.ReactNode = tooltip ? (
    <Tooltip placement="bottom" title={tooltip}>
      {chartTitle}
    </Tooltip>
  ) : (
    chartTitle
  );

  return (
    <div
      className={cx(
        "visualization",
        stale && "visualization-faded",
        preCalcGraphSize && "visualization-graph-sizing",
      )}
    >
      <div className={cx("visualization-header")}>{tooltipNode}</div>
      <div
        className={cx(
          "visualization-content",
          loading && "visualization-loading",
        )}
      >
        {loading ? (
          <img className={cx("visualization-spinner")} src={spinner} />
        ) : (
          children
        )}
      </div>
    </div>
  );
};
