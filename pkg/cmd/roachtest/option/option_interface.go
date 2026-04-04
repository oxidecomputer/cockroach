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

package option

// Option is an interface used to configure operations against a Cluster.
//
// TODO(tschottdorf): Consider using a more idiomatic approach in which options
// act upon a config struct:
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option interface {
	Option() // TODO(tbg): lower-case this again when all opts are in this pkg
}
