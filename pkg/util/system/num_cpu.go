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

package system

import "runtime"

// NumCPU returns the number of logical CPUs usable by the current process.
//
// !!Note!! If you are considering using this to scale parallelism with the
// machine size, use runtime.GOMAXPROCS(0) instead. The latter is better because
// GOMAXPROCS is reduced with certain test runs (with the race detector); it can
// also be reduced in containerized environments.
//
// The set of available CPUs is checked by querying the operating system
// at process startup. Changes to operating system CPU allocation after
// process startup are not reflected.
func NumCPU() int {
	return runtime.NumCPU()
}
