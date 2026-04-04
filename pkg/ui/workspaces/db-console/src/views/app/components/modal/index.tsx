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

import { Modal, Button } from "antd";
import React, { Fragment } from "react";
import "./styles.styl";
import { ModalProps } from "antd/lib/modal";

interface ICustomModalProps extends ModalProps {
  children?: React.ReactNode;
  trigger?: React.ReactChildren | React.ReactNode;
  triggerStyle?: string;
  triggerTitle?: string;
}

interface ICustomModalState {
  visible: boolean;
}

class CustomModal extends React.Component<
  ICustomModalProps,
  ICustomModalState
> {
  state = { visible: false };

  showModal = () => {
    this.setState({
      visible: true,
    });
  };

  handleOk = () => {
    this.setState({
      visible: false,
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
    });
  };

  render() {
    const {
      trigger,
      visible,
      children,
      triggerStyle,
      triggerTitle,
    } = this.props;
    return (
      <Fragment>
        {trigger ? (
          trigger
        ) : (
          <a onClick={this.showModal} className={`${triggerStyle}`}>
            {triggerTitle}
          </a>
        )}
        <Modal
          visible={trigger ? visible : this.state.visible}
          onOk={this.handleOk}
          onCancel={this.handleCancel}
          className="custom--modal"
          maskStyle={{
            background: "rgba(71, 88, 114, 0.73)",
          }}
          footer={
            <Button
              type="link"
              className="custom--modal__close--button"
              onClick={this.handleCancel}
            >
              Done
            </Button>
          }
          {...this.props}
        >
          {children}
        </Modal>
      </Fragment>
    );
  }
}

export default CustomModal;
