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

package reduce

import (
	"log"
	"os"
	"testing"

	"github.com/cockroachdb/datadriven"
)

// Walk walks path for datadriven files and calls RunTest on them.
func Walk(
	t *testing.T,
	path string,
	filter func(string) (string, error),
	interesting func(contains string) InterestingFn,
	mode Mode,
	cr ChunkReducer,
	passes []Pass,
) {
	datadriven.Walk(t, path, func(t *testing.T, path string) {
		RunTest(t, path, filter, interesting, mode, cr, passes)
	})
}

// RunTest executes reducer commands. If the optional filter func is not nil,
// it is invoked on all "reduce" commands before the reduction algorithm is
// started. In verbose mode, output is written to stderr. Supported commands:
//
// "contains": Sets the substring that must be contained by an error message
// for a test case to be marked as interesting. This variable is passed to the
// `interesting` func.
//
// "reduce": Creates an InterestingFn based on the current value of contains
// and executes the reducer. Outputs the reduced case.
func RunTest(
	t *testing.T,
	path string,
	filter func(string) (string, error),
	interesting func(contains string) InterestingFn,
	mode Mode,
	cr ChunkReducer,
	passes []Pass,
) {
	var contains string
	var logger *log.Logger
	if testing.Verbose() {
		logger = log.New(os.Stderr, "", 0)
	}
	datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
		switch d.Cmd {
		case "contains":
			contains = d.Input
			return ""
		case "reduce":
			input := d.Input
			if filter != nil {
				var err error
				input, err = filter(input)
				if err != nil {
					t.Fatal(err)
				}
			}
			output, err := Reduce(logger, input, interesting(contains), 0, mode, cr, passes...)
			if err != nil {
				t.Fatal(err)
			}
			return output + "\n"
		default:
			t.Fatalf("unknown command: %s", d.Cmd)
			return ""
		}
	})
}
