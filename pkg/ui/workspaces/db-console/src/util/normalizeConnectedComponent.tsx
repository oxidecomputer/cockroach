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

import React, { ExoticComponent } from "react";

/*
 * normalizeConnectedComponent function returns react element created by wrapping Connected component (which in fact is
 * not a 'valid' react component (see: ) and provided properties.
 * It is required for passing correct components to Route component.
 * For more details see: @types/react/index.d.ts:314
 * > "However, we have no way of telling the JSX parser that it's a JSX element type or its props other than
 * > by pretending to be a normal component."
 * */
export const normalizeConnectedComponent = (
  ConnectedComponent: ExoticComponent,
) => (props: React.ComponentProps<any>) => <ConnectedComponent {...props} />;
