// Copyright 2020 The Cockroach Authors.
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
import { Link, RouteComponentProps, withRouter } from "react-router-dom";

import { Text, TextTypes } from "src/components";
import { trustIcon } from "src/util/trust";
import ArrowLeftIcon from "!!raw-loader!assets/arrowLeft.svg";
import "./headerSection.styl";

export interface HeaderSectionProps {
  title: string;
  navigationBackConfig?: {
    text: string;
    path: string;
  };
}

const HeaderSection: React.FC<HeaderSectionProps &
  RouteComponentProps> = props => {
  const { navigationBackConfig, title } = props;
  return (
    <div className="header-section">
      {navigationBackConfig && (
        <div className="header-section__back-link">
          <span
            className="header-section__back-icon"
            dangerouslySetInnerHTML={trustIcon(ArrowLeftIcon)}
          />
          <Link to={navigationBackConfig.path}>
            {navigationBackConfig.text}
          </Link>
        </div>
      )}
      <div className="header-section__title">
        <Text textType={TextTypes.Heading3}>{title}</Text>
      </div>
    </div>
  );
};

export default withRouter(HeaderSection);
