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
import { TextInput, Button } from "src/components";
import { isValidEmail } from "src/util/validation/isValidEmail";

import "./emailSubscriptionForm.styl";

interface EmailSubscriptionFormState {
  emailAddress: string | undefined;
  canSubmit: boolean;
}

interface EmailSubscriptionFormProps {
  onSubmit?: (emailAddress: string) => void;
}

export class EmailSubscriptionForm extends React.Component<
  EmailSubscriptionFormProps,
  EmailSubscriptionFormState
> {
  constructor(props: EmailSubscriptionFormProps) {
    super(props);
    this.state = {
      emailAddress: undefined,
      canSubmit: false,
    };
  }

  handleSubmit = () => {
    if (this.state.canSubmit) {
      this.props.onSubmit(this.state.emailAddress);
      this.setState({
        emailAddress: "",
        canSubmit: false,
      });
    }
  };

  handleChange = (value: string) => {
    this.handleEmailValidation(value);
    this.setState({
      emailAddress: value,
    });
  };

  handleEmailValidation = (value: string) => {
    const isCorrectEmail = isValidEmail(value);
    const isEmpty = value.length === 0;

    this.setState({
      canSubmit: isCorrectEmail && !isEmpty,
    });

    if (isCorrectEmail || isEmpty) {
      return undefined;
    }
    return "Invalid email address.";
  };

  render() {
    const { canSubmit, emailAddress } = this.state;
    return (
      <div className="email-subscription-form">
        <TextInput
          name="email"
          className="email-subscription-form__input"
          placeholder="Enter your email"
          validate={this.handleEmailValidation}
          onChange={this.handleChange}
          value={emailAddress}
        />
        <Button
          type={"primary"}
          onClick={this.handleSubmit}
          disabled={!canSubmit}
          className="email-subscription-form__submit-button"
        >
          Sign up
        </Button>
      </div>
    );
  }
}
