// Copyright 2019 The Cockroach Authors.
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

import { Checkbox, Divider } from "antd";
import Dropdown, { DropdownOption } from "src/views/shared/components/dropdown";
import React from "react";
import { RouteComponentProps, withRouter } from "react-router-dom";

import { getMatchParamByName } from "src/util/query";
import { NetworkFilter, NetworkSort } from "..";
import { Filter } from "../filter";
import "./sort.styl";

interface ISortProps {
  onChangeFilter: (key: string, value: string) => void;
  onChangeCollapse: (checked: boolean) => void;
  deselectFilterByKey: (key: string) => void;
  collapsed: boolean;
  sort: NetworkSort[];
  filter: NetworkFilter;
}

class Sort extends React.Component<ISortProps & RouteComponentProps, {}> {
  onChange = ({ target }: any) => this.props.onChangeCollapse(target.checked);

  pageView = () => {
    const { match } = this.props;
    const nodeId = getMatchParamByName(match, "node_id");
    return nodeId || "cluster";
  };

  navigateTo = (selected: DropdownOption) => {
    this.props.onChangeCollapse(false);
    this.props.history.push(`/reports/network/${selected.value}`);
  };

  componentDidMount() {
    this.setDefaultSortValue("region");
  }

  setDefaultSortValue = (sortValue: string) => {
    const isDefaultValuePresent = this.getSortValues(this.props.sort).find(
      e => e.value === sortValue,
    );
    if (isDefaultValuePresent) {
      this.navigateTo(isDefaultValuePresent);
    }
  };

  getSortValues = (sort: NetworkSort[]) =>
    sort.map(value => {
      return {
        value: value.id,
        label: value.id.replace(/^[a-z]/, m => m.toUpperCase()),
      };
    });

  render() {
    const {
      collapsed,
      sort,
      onChangeFilter,
      deselectFilterByKey,
      filter,
      match,
    } = this.props;
    const nodeId = getMatchParamByName(match, "node_id");
    return (
      <div className="Sort-latency">
        <Dropdown
          title="Sort By"
          options={this.getSortValues(sort)}
          selected={this.pageView()}
          onChange={this.navigateTo}
          className="Sort-latency__dropdown--spacing"
        />
        <Filter
          sort={sort}
          onChangeFilter={onChangeFilter}
          deselectFilterByKey={deselectFilterByKey}
          filter={filter}
          dropDownClassName="Sort-latency__dropdown--spacing"
        />
        <Divider type="vertical" style={{ height: "100%" }} />
        <Checkbox
          disabled={!nodeId || nodeId === "cluster"}
          checked={collapsed}
          onChange={this.onChange}
        >
          Collapse Nodes
        </Checkbox>
      </div>
    );
  }
}

export default withRouter(Sort);
