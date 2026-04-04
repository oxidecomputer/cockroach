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

package log

import (
	"strings"
	"testing"

	"github.com/cockroachdb/datadriven"
)

func TestReadLogFormat(t *testing.T) {
	datadriven.RunTest(t, "testdata/read_header",
		func(t *testing.T, td *datadriven.TestData) string {
			switch td.Cmd {
			case "log":
				_, format, err := ReadFormatFromLogFile(strings.NewReader(td.Input))
				if err != nil {
					td.Fatalf(t, "error while reading format from the log file: %v", err)
				}
				return format
			default:
				t.Fatalf("unknown directive: %q", td.Cmd)
			}
			// unreachable
			return ""
		})
}
