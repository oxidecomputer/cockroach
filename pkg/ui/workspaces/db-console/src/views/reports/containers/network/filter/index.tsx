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

import { Checkbox, Select } from "antd";
import Dropdown, { arrowRenderer } from "src/views/shared/components/dropdown";
import React from "react";
import classNames from "classnames";
import { NetworkFilter, NetworkSort } from "..";
import "./filter.styl";

interface IFilterProps {
  onChangeFilter: (key: string, value: string) => void;
  deselectFilterByKey: (key: string) => void;
  sort: NetworkSort[];
  filter: NetworkFilter;
  dropDownClassName?: string;
}

interface IFilterState {
  opened: boolean;
  width: number;
}

export class Filter extends React.Component<IFilterProps, IFilterState> {
  state = {
    opened: false,
    width: window.innerWidth,
  };

  private rangeContainer = React.createRef<HTMLDivElement>();

  componentDidMount() {
    window.addEventListener("resize", this.updateDimensions);
  }

  componentWillUnmount() {
    window.removeEventListener("resize", this.updateDimensions);
  }

  updateDimensions = () => {
    this.setState({
      width: window.innerWidth,
    });
  };

  onChange = (key: string, value: string) => () =>
    this.props.onChangeFilter(key, value);

  onDeselect = (key: string) => () => this.props.deselectFilterByKey(key);

  renderSelectValue = (id: string) => {
    const { filter } = this.props;

    if (filter && filter[id]) {
      const value = (key: string) =>
        `${filter[id].length} ${this.firstLetterToUpperCase(key)} Selected`;
      switch (true) {
        case filter[id].length === 1 && id === "cluster":
          return value("Node");
        case filter[id].length === 1:
          return value(id);
        case filter[id].length > 1 && id === "cluster":
          return value("Nodes");
        case filter[id].length > 1:
          return value(`${id}s`);
        default:
          return;
      }
    }
    return;
  };

  firstLetterToUpperCase = (value: string) =>
    value.replace(/^[a-z]/, m => m.toUpperCase());

  renderSelect = () => {
    const { sort, filter } = this.props;
    return sort.map(value => (
      <div style={{ width: "100%" }} className="select__container">
        <p className="filter--label">{`${
          value.id === "cluster"
            ? "Nodes"
            : this.firstLetterToUpperCase(value.id)
        }`}</p>
        <Select
          style={{ width: "100%" }}
          placeholder={`Filter ${
            value.id === "cluster" ? "node" : value.id
          }(s)`}
          value={this.renderSelectValue(value.id)}
          dropdownRender={_ => (
            <div onMouseDown={e => e.preventDefault()}>
              <div className="select-selection__deselect">
                <a onClick={this.onDeselect(value.id)}>Deselect all</a>
              </div>
              {value.filters.map(val => {
                const checked =
                  filter &&
                  filter[value.id] &&
                  filter[value.id].indexOf(val.name) !== -1;
                return (
                  <div className="filter__checkbox">
                    <Checkbox
                      checked={checked}
                      onChange={this.onChange(value.id, val.name)}
                    />
                    <a
                      className={`filter__checkbox--label ${
                        checked ? "filter__checkbox--label__active" : ""
                      }`}
                      onClick={this.onChange(value.id, val.name)}
                    >{`${value.id === "cluster" ? "N" : ""}${val.name}: ${
                      val.address
                    }`}</a>
                  </div>
                );
              })}
            </div>
          )}
        />
      </div>
    ));
  };

  render() {
    const { opened, width } = this.state;
    const { dropDownClassName } = this.props;
    const containerLeft = this.rangeContainer.current
      ? this.rangeContainer.current.getBoundingClientRect().left
      : 0;
    const left =
      width >= containerLeft + 240 ? 0 : width - (containerLeft + 240);
    return (
      <div className="Filter-latency">
        <Dropdown
          title="Filter"
          options={[]}
          selected=""
          className={classNames(
            { dropdown__focused: opened },
            dropDownClassName,
          )}
          content={
            <div ref={this.rangeContainer} className="Range">
              <div
                className="click-zone"
                onClick={() => this.setState({ opened: !opened })}
              />
              {opened && (
                <div
                  className="trigger-container"
                  onClick={() => this.setState({ opened: false })}
                />
              )}
              <div className="trigger-wrapper">
                <div
                  className={`trigger Select ${(opened && "is-open") || ""}`}
                >
                  <div className="Select-control">
                    <div className="Select-arrow-zone">
                      {arrowRenderer({ isOpen: opened })}
                    </div>
                  </div>
                </div>
                {opened && (
                  <div className="multiple-filter__selection" style={{ left }}>
                    {this.renderSelect()}
                  </div>
                )}
              </div>
            </div>
          }
        />
      </div>
    );
  }
}
