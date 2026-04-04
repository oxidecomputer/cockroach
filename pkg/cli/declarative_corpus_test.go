// Copyright 2022 The Cockroach Authors.
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

package cli

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/datadriven"
)

// This test doctoring a secure cluster.
func TestDeclarativeCorpus(t *testing.T) {
	defer leaktest.AfterTest(t)()
	c := NewCLITest(TestCLIParams{T: t, NoServer: true})
	defer c.Cleanup()

	t.Run("declarative corpus validation standalone command", func(t *testing.T) {
		out, err := c.RunWithCapture("debug declarative-corpus-validate testdata/declarative-corpus/corpus")
		if err != nil {
			t.Fatal(err)
		}

		// Using datadriven allows TESTFLAGS=-rewrite.
		datadriven.RunTest(t, testutils.TestDataPath(t, "declarative-corpus", "corpus_expected"), func(t *testing.T, td *datadriven.TestData) string {
			return out
		})
	})
}
