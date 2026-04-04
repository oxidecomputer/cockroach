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

package logconfig

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cockroachdb/datadriven"
	"gopkg.in/yaml.v2"
)

func TestExport(t *testing.T) {
	datadriven.RunTest(t, "testdata/export", func(t *testing.T, d *datadriven.TestData) string {
		var onlyChans ChannelList
		if d.HasArg("only-channels") {
			var s string
			d.ScanArgs(t, "only-channels", &s)
			chs, err := parseChannelList(s)
			if err != nil {
				t.Fatal(err)
			}
			onlyChans.Channels = chs
		}

		c := DefaultConfig()
		if err := yaml.UnmarshalStrict([]byte(d.Input), &c); err != nil {
			t.Fatal(err)
		}
		defaultDir := "/default-dir"
		var buf bytes.Buffer
		if err := c.Validate(&defaultDir); err != nil {
			t.Fatal(err)
		} else {
			uml, key := c.Export(onlyChans)
			buf.WriteString(uml)
			fmt.Fprintf(&buf, "# http://www.plantuml.com/plantuml/uml/%s\n", key)
		}
		return buf.String()
	})
}
