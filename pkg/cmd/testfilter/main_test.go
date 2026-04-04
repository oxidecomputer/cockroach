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

package main

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/datadriven"
)

func TestFilter(t *testing.T) {
	defer leaktest.AfterTest(t)()

	datadriven.Walk(t, testutils.TestDataPath(t), func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			in := strings.NewReader(td.Input)
			var out strings.Builder
			var mode modeT
			if err := mode.Set(td.Cmd); err != nil {
				return err.Error()
			}
			if err := filter(in, &out, mode); err != nil {
				return err.Error()
			}
			// At the time of writing, datadriven garbles the test files when
			// rewriting a "\n" output, so make sure we never have trailing
			// newlines.
			return strings.TrimRight(out.String(), "\r\n")
		})
	})
}
