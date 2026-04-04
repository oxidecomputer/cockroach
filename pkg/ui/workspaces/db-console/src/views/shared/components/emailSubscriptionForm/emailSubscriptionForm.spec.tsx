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
import { assert } from "chai";
import { mount, ReactWrapper } from "enzyme";
import sinon, { SinonSpy } from "sinon";
import classNames from "classnames/bind";

import "src/enzymeInit";
import { EmailSubscriptionForm } from "./index";
import buttonStyles from "src/components/button/button.module.styl";

const cx = classNames.bind(buttonStyles);
const sandbox = sinon.createSandbox();

describe("EmailSubscriptionForm", () => {
  let wrapper: ReactWrapper;
  let onSubmitHandler: SinonSpy;

  beforeEach(() => {
    sandbox.reset();
    onSubmitHandler = sandbox.spy();
    wrapper = mount(<EmailSubscriptionForm onSubmit={onSubmitHandler} />);
  });

  describe("when correct email", () => {
    it("provides entered email on submit callback", () => {
      const emailAddress = "foo@bar.com";
      const inputComponent = wrapper.find("input.crl-input__text").first();
      inputComponent.simulate("change", { target: { value: emailAddress } });
      const buttonComponent = wrapper
        .find(`button.${cx("crl-button")}`)
        .first();
      buttonComponent.simulate("click");

      onSubmitHandler.calledOnceWith(emailAddress);
    });
  });

  describe("when invalid email", () => {
    beforeEach(() => {
      const emailAddress = "foo";
      const inputComponent = wrapper.find("input.crl-input__text").first();
      inputComponent.simulate("change", { target: { value: emailAddress } });
      inputComponent.simulate("blur");
    });

    it("doesn't call onSubmit callback", () => {
      const buttonComponent = wrapper
        .find(`button.${cx("crl-button")}`)
        .first();
      buttonComponent.simulate("click");
      assert.isTrue(onSubmitHandler.notCalled);
    });

    it("submit button is disabled", () => {
      const buttonComponent = wrapper
        .find(`button.${cx("crl-button")}.${cx("crl-button--disabled")}`)
        .first();
      assert.isTrue(buttonComponent.exists());
    });

    it("validation message is shown", () => {
      const validationMessageWrapper = wrapper
        .find(".crl-input__text--error-message")
        .first();
      assert.isTrue(validationMessageWrapper.exists());
      assert.equal(validationMessageWrapper.text(), "Invalid email address.");
    });
  });
});
