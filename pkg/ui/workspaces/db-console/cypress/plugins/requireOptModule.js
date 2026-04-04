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

const OPT_NODE_MODULES_PATH = "./opt/node_modules";
const path = require('path');

/**
 * @description Resolves module loading from custom `node_modules` location.
 * It is required only for .js files which aren't processed by Webpack.
 * @param module {string} module name to load from ./opt/node_modules
 * @return {*} required module
 */
function requireOptModule(module) {
  return require(path.resolve(OPT_NODE_MODULES_PATH, module));
}

module.exports = requireOptModule;
module.exports.OPT_NODE_MODULES_PATH = OPT_NODE_MODULES_PATH;
