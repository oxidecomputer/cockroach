// Copyright 2022 The Cockroach Authors.
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

import moment from "moment";
import React, { useState, useEffect } from "react";

type Props = {
  children: React.ReactElement;
  delay?: moment.Duration;
};

export const Delayed = ({
  children,
  delay = moment.duration(10, "s"),
}: Props): React.ReactElement => {
  const [isShown, setIsShown] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsShown(true);
    }, delay.asMilliseconds());
    return () => clearTimeout(timer);
  }, [delay]);

  return isShown ? children : null;
};
