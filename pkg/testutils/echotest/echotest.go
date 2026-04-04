// Copyright 2022 The Cockroach Authors.
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

package echotest

import (
	"testing"

	"github.com/cockroachdb/datadriven"
)

// Require checks that the string matches what is found in the file located at
// the provided path. The file must follow the datadriven format:
//
// echo
// ----
// <output of exp>
//
// The contents of the file can be updated automatically using datadriven's
// -rewrite flag.
func Require(t *testing.T, act, path string) {
	var ran bool
	datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
		if d.Cmd != "echo" {
			return "only 'echo' is supported"
		}
		ran = true
		return act
	})
	if !ran {
		// Guard against a possible error in which the file is created, then datadriven
		// is invoked with -rewrite to seed it (which it does not do, since there is
		// no directive in the file), and then also the tests pass despite not checking
		// anything.
		t.Errorf("no tests run for %s, is the file empty?", path)
	}
}
