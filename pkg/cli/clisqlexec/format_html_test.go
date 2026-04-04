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

package clisqlexec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestRenderHTML(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	cols := []string{"colname"}
	align := "d"
	rows := [][]string{
		{"<b>foo</b>"},
		{"bar"},
	}

	type testCase struct {
		reporter htmlReporter
		out      string
	}

	testCases := []testCase{
		{
			reporter: htmlReporter{},
			out: `<table>
<thead><tr><th>colname</th></tr></thead>
<tbody>
<tr><td><b>foo</b></td></tr>
<tr><td>bar</td></tr>
</tbody>
</table>
`,
		},
		{
			reporter: htmlReporter{escape: true},
			out: `<table>
<thead><tr><th>colname</th></tr></thead>
<tbody>
<tr><td>&lt;b&gt;foo&lt;/b&gt;</td></tr>
<tr><td>bar</td></tr>
</tbody>
</table>
`,
		},
		{
			reporter: htmlReporter{rowStats: true},
			out: `<table>
<thead><tr><th>row</th><th>colname</th></tr></thead>
<tbody>
<tr><td>1</td><td><b>foo</b></td></tr>
<tr><td>2</td><td>bar</td></tr>
</tbody>
<tfoot><tr><td colspan=2>2 rows</td></tr></tfoot></table>
`,
		},
		{
			reporter: htmlReporter{escape: true, rowStats: true},
			out: `<table>
<thead><tr><th>row</th><th>colname</th></tr></thead>
<tbody>
<tr><td>1</td><td>&lt;b&gt;foo&lt;/b&gt;</td></tr>
<tr><td>2</td><td>bar</td></tr>
</tbody>
<tfoot><tr><td colspan=2>2 rows</td></tr></tfoot></table>
`,
		},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("escape=%v/rowStats=%v", tc.reporter.escape, tc.reporter.rowStats)
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			err := render(&tc.reporter, &buf, ioutil.Discard,
				cols, NewRowSliceIter(rows, align),
				nil /* completedHook */, nil /* noRowsHook */)
			if err != nil {
				t.Fatal(err)
			}
			if tc.out != buf.String() {
				t.Errorf("expected:\n%s\ngot:\n%s", tc.out, buf.String())
			}
		})
	}
}
