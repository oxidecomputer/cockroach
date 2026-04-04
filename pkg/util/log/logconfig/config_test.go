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
	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	datadriven.RunTest(t, "testdata/yaml", func(t *testing.T, d *datadriven.TestData) string {
		var c Config
		if err := yaml.UnmarshalStrict([]byte(d.Input), &c); err != nil {
			return fmt.Sprintf("ERROR: %v\n", err)
		}
		t.Logf("%# v", pretty.Formatter(c))
		var buf bytes.Buffer
		b, err := yaml.Marshal(&c)
		if err != nil {
			fmt.Fprintf(&buf, "ERROR: %v\n", err)
		} else {
			fmt.Fprintf(&buf, "%s", string(b))
		}
		return buf.String()
	})
}
