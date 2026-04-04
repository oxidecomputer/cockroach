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

// Package clisqlcfg defines configuration settings and mechanisms for
// instances of the SQL shell.
//
// This package is intended to be used as follows:
//
// 1. instantiate a configuration with `NewDefaultConfig()`.
//
// 2. load customizations from e.g. command-line flags, env vars, etc.
//
// 3. validate the configuration and open the input/output streams via
//    `(*Context).Open()`. Defer a call to the returned cleanup function.
//
// 4. open a client connection via `(*Context).MakeConn()`.
//    Note: this must occur after the call to `Open()`, as the configuration
//    may not be ready before that point.
//
// 5. call `(*Context).Run()`.
//
package clisqlcfg
