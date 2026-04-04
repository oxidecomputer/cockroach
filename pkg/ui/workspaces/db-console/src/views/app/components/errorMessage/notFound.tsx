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
import Helmet from "react-helmet";
import "./errorMessage.styl";
import NotFoundImg from "assets/not-found.svg";

function NotFound() {
  return (
    <main className="error-message-page">
      <Helmet title="Not Found" />
      <div className="error-message-page__content">
        <img
          className="error-message-page__img"
          src={NotFoundImg}
          alt="404 Error"
        />
        <div className="error-message-page__body">
          <div className="error-message-page__message">Whoops!</div>
          <p>
            We can&apos;t find the page you are looking for. You may have typed
            the wrong address or found a broken link.
          </p>
        </div>
      </div>
    </main>
  );
}

export default NotFound;
