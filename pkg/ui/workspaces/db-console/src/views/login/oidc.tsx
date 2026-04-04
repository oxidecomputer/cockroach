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

import { LoginAPIState } from "oss/src/redux/login";
import { Button } from "src/components";
import { RouteComponentProps, withRouter } from "react-router-dom";

const OIDC_LOGIN_PATH = "/oidc/v1/login";

const OIDCLoginButton = ({ loginState }: { loginState: LoginAPIState }) => {
  return (
    <a href={OIDC_LOGIN_PATH}>
      <Button
        type="secondary"
        className="submit-button-oidc"
        disabled={loginState.inProgress}
        textAlign={"center"}
      >
        {loginState.oidcButtonText}
      </Button>
    </a>
  );
};

const OIDCLogin: React.FC<{
  loginState: LoginAPIState;
} & RouteComponentProps> = props => {
  const oidcAutoLoginQuery = new URLSearchParams(props.location.search).get(
    "oidc_auto_login",
  );
  if (props.loginState.oidcLoginEnabled) {
    if (props.loginState.oidcAutoLogin && !(oidcAutoLoginQuery === "false")) {
      window.location.replace(OIDC_LOGIN_PATH);
    }
    return <OIDCLoginButton loginState={props.loginState} />;
  }
  return null;
};

export const OIDCLoginConnected = withRouter(OIDCLogin);
