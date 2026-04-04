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

package install

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/datadriven"
	"github.com/stretchr/testify/require"
)

func TestExecStartTemplate(t *testing.T) {
	startData1 := startTemplateData{
		LogDir: "./path with spaces/logs/$THIS_DOES_NOT_EVER_GET_EXPANDED",
		KeyCmd: `echo foo && \
echo bar $HOME`,
		EnvVars:   []string{"ROACHPROD=1/tigtag", "COCKROACH=foo", "ROCKCOACH=17%"},
		Binary:    "./cockroach",
		Args:      []string{`start`, `--log`, `file-defaults: {dir: '/path with spaces/logs', exit-on-error: false}`},
		MemoryMax: "81%",
		Local:     true,
	}
	datadriven.Walk(t, testutils.TestDataPath(t, "start"), func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			if td.Cmd != "data1" {
				t.Fatalf("unsupported")
			}
			out, err := execStartTemplate(startData1)
			require.NoError(t, err)
			return out
		})
	})
}
