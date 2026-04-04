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

import React, { ErrorInfo } from "react";
import Helmet from "react-helmet";
import "./errorMessage.styl";
import SleepyMoonImg from "assets/sleepy-moon.svg";

interface ErrorBoundaryProps {
  onCatch?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | undefined;
}

// ErrorBoundary with image and text message.
export default class ErrorBoundary extends React.Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: undefined,
    };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Console.error for developer visibility.
    console.error(error);
    this.props.onCatch && this.props.onCatch(error, errorInfo);
  }

  render() {
    if (!this.state.hasError) {
      return this.props.children;
    }
    return (
      <main className="error-message-page">
        <Helmet title="Error" />
        <div className="error-message-page__content">
          <img className="error-message-page__img" src={SleepyMoonImg} />
          <div className="error-message-page__body">
            <div className="error-message-page__message">
              Something went wrong.
            </div>
            <p>
              There is a problem loading the component of this page. Try
              refreshing the page.
            </p>
          </div>
        </div>
      </main>
    );
  }
}
