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
import { RouteComponentProps, withRouter } from "react-router-dom";
import { connect } from "react-redux";

import { AdminUIState } from "src/redux/state";
import { selectLoginState, LoginState, getLoginPage } from "src/redux/login";

interface RequireLoginProps {
  loginState: LoginState;
}

class RequireLogin extends React.Component<
  RouteComponentProps & RequireLoginProps
> {
  componentDidMount() {
    this.checkLogin();
  }

  componentDidUpdate() {
    this.checkLogin();
  }

  checkLogin() {
    const { location, history } = this.props;

    if (!this.hideLoginPage()) {
      history.push(getLoginPage(location));
    }
  }

  hideLoginPage() {
    return this.props.loginState.hideLoginPage();
  }

  render() {
    if (!this.hideLoginPage()) {
      return null;
    }

    return this.props.children;
  }
}

const RequireLoginConnected = withRouter(
  connect((state: AdminUIState) => {
    return {
      loginState: selectLoginState(state),
    };
  })(RequireLogin),
);

export default RequireLoginConnected;
