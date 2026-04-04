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

package schemaexpr

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/datadriven"
)

func TestParseComputedColumnRewrites(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	path := testutils.TestDataPath(t, "computed_column_rewrites")
	datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
		switch d.Cmd {
		case "parse":
			rewrites, err := ParseComputedColumnRewrites(d.Input)
			if err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			var keys []string
			for k := range rewrites {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			var buf bytes.Buffer
			for _, k := range keys {
				fmt.Fprintf(&buf, "(%v) -> (%v)\n", k, tree.Serialize(rewrites[k]))
			}
			return buf.String()

		default:
			t.Fatalf("unsupported command %s", d.Cmd)
			return ""
		}
	})
}
