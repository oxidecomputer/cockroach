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

package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/datadriven"
)

func TestPretty(t *testing.T) {
	datadriven.Walk(t, testutils.TestDataPath(t), func(t *testing.T, path string) {
		datadriven.RunTest(t, path, prettyTest)
	})
}

func prettyTest(t *testing.T, d *datadriven.TestData) string {
	switch d.Cmd {
	case "pretty":
		n := defaultWidth
		if d.HasArg("n") {
			d.ScanArgs(t, "n", &n)
		}
		exprgen := d.HasArg("expr")
		s, err := prettyify(strings.NewReader(d.Input), n, exprgen)
		if err != nil {
			return fmt.Sprintf("ERROR: %s", err)
		}

		// Verify we round trip correctly by ensuring non-whitespace
		// scanner tokens are encountered in the same order.
		{
			origToks := toTokens(d.Input)
			prettyToks := toTokens(s)
			for i, tok := range origToks {
				if i >= len(prettyToks) {
					t.Fatalf("pretty ended early after %d tokens", i+1)
				}
				if prettyToks[i] != tok {
					t.Log(s)
					t.Logf("expected %q", tok)
					t.Logf("got %q", prettyToks[i])
					t.Fatalf("token %d didn't match", i+1)
				}
			}
			if len(prettyToks) > len(origToks) {
				t.Fatalf("orig ended early after %d tokens", len(origToks))
			}
		}
		// Verify lines aren't too long.
		{
			for i, line := range strings.Split(s, "\n") {
				if strings.HasPrefix(line, "#") {
					continue
				}
				if len(line) > defaultWidth {
					t.Errorf("line %d is %d chars, expected <= %d:\n%s", i+1, len(line), defaultWidth, line)
				}
			}
		}

		return s
	default:
		t.Fatal("unknown command")
		return ""
	}
}
