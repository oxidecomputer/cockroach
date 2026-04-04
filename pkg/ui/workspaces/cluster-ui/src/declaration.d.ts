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

declare module "*.module.scss" {
  const content: { [className: string]: string };
  export default content;
}
declare module "*.png";
declare module "*.gif";
declare module "*.scss";
declare module "*.svg";

type ConstructorType = new (...args: any) => any;

type FirstConstructorParameter<
  P extends ConstructorType
> = ConstructorParameters<P>[0];

type Tuple<T> = [T, T];

type Dictionary<V> = {
  [key: string]: V;
};

declare module "highlight.js/lib/core";
declare module "highlight.js/lib/languages/pgsql";
