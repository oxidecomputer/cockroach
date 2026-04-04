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

//go:build crdb_test && !crdb_test_off
// +build crdb_test,!crdb_test_off

package buildutil

// CrdbTestBuild is a flag that is set to true if the binary was compiled
// with the 'crdb_test' build tag (which is the case for all test targets). This
// flag can be used to enable expensive checks, test randomizations, or other
// metamorphic-style perturbations that will not affect test results but will
// exercise different parts of the code.
const CrdbTestBuild = true
