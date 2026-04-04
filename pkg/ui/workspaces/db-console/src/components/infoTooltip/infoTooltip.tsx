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

import * as React from "react";

import { ToolTipWrapper } from "src/views/shared/components/toolTip";

import "./infoTooltip.styl";

export const InfoTooltip = (props: { text: React.ReactNode }) => {
  const { text } = props;
  return (
    <div className="info-tooltip__tooltip">
      <ToolTipWrapper text={text}>
        <div className="info-tooltip__tooltip-hover-area">
          <div className="info-tooltip__info-icon">i</div>
        </div>
      </ToolTipWrapper>
    </div>
  );
};
