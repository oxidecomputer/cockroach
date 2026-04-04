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

package execgen

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/datadriven"
	"github.com/dave/dst/decorator"
)

// Walk walks path for datadriven files and calls RunTest on them.
func TestExecgen(t *testing.T) {
	datadriven.Walk(t, testutils.TestDataPath(t), func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			f, err := decorator.Parse(d.Input)
			if err != nil {
				t.Fatal(err)
			}
			switch d.Cmd {
			case "inline":
				inlineFuncs(f)
			case "template":
				expandTemplates(f)
			default:
				t.Fatalf("unknown command: %s", d.Cmd)
				return ""
			}
			var sb strings.Builder
			if err := decorator.Fprint(&sb, f); err != nil {
				return err.Error()
			}
			return sb.String()
		})
	})
}
