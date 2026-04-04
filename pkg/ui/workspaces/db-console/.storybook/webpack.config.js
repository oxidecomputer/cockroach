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

const custom = require("../webpack.app.js");

const appConfig = custom({dist: "ccl"}, {mode: "development"});

module.exports = async ({ config, mode }) => {
  return {
    ...config,
    resolve: {
      ...config.resolve,
      modules: appConfig.resolve.modules,
      extensions: appConfig.resolve.extensions,
    },
    module: {
      rules: [
        ...appConfig.module.rules,
      ]
    },
    plugins: [
      ...config.plugins,
      // Import 'vendors' library only
      appConfig.plugins[2]
    ]
  };
};
