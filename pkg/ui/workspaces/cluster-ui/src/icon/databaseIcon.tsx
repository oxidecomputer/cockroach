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

export const DatabaseIcon = ({ className, ...props }: IconProps) => (
  <svg viewBox="0 0 14 14" className={className} {...props}>
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M12.25 1.167H1.75a.583.583 0 00-.583.583v10.5c0 .322.26.583.583.583h10.5a.583.583 0 00.583-.583V1.75a.583.583 0 00-.583-.583zM1.75 0A1.75 1.75 0 000 1.75v10.5C0 13.216.784 14 1.75 14h10.5A1.75 1.75 0 0014 12.25V1.75A1.75 1.75 0 0012.25 0H1.75z"
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M3.662 13.417V1h1.239v3.292H13V5.49H4.9V8.5H13v1.25H4.9v3.667H3.663z"
    />
  </svg>
);
