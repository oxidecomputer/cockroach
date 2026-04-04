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

import * as React from "react";

export interface IconProps {
  fill?: string;
}

export function CaretDown(props: IconProps) {
  const { fill } = props;

  return (
    <svg width={8} height={6} viewBox="0 0 8 6" fill="none" {...props}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M1.39.45a.667.667 0 10-1.003.878l3.111 3.555a.667.667 0 001.004 0l3.11-3.555A.667.667 0 106.61.45L4 3.432 1.39.45z"
        fill={fill}
      />
    </svg>
  );
}

CaretDown.defaultProps = {
  fill: "#475872",
};
