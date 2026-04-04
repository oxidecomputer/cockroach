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

import React from "react";
import moment from "moment";
import { ToolTipWrapper } from "src/views/shared/components/toolTip";
import { DATE_FORMAT } from "src/util/format";
import { google } from "src/js/protos";
import ITimestamp = google.protobuf.ITimestamp;

interface HighwaterProps {
  timestamp: ITimestamp;
  decimalString: string;
}

export class HighwaterTimestamp extends React.PureComponent<HighwaterProps> {
  render() {
    if (!this.props.timestamp) {
      return null;
    }
    let highwaterMoment = moment(
      this.props.timestamp.seconds.toNumber() * 1000,
    );
    // It's possible due to client clock skew that this timestamp could be in
    // the future. To avoid confusion, set a maximum bound of now.
    const now = moment();
    if (highwaterMoment.isAfter(now)) {
      highwaterMoment = now;
    }

    return (
      <ToolTipWrapper text={highwaterMoment.format(DATE_FORMAT)}>
        {this.props.decimalString}
      </ToolTipWrapper>
    );
  }
}
