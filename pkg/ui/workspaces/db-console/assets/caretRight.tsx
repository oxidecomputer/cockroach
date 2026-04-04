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

import * as React from "react";

interface IconProps {
  className?: string;
}

export const CaretRight = ({ className, ...props }: IconProps) => (
  <svg viewBox="0 0 11 17" className={className} {...props}>
    <path
      fillRule="evenodd"
      d="M.512 14.371a1.5 1.5 0 1 0 1.976 2.258l8-7a1.5 1.5 0 0 0 0-2.258l-8-7A1.5 1.5 0 0 0 .512 2.63L7.222 8.5l-6.71 5.871z"
      clipRule="evenodd"
    />
  </svg>
);
