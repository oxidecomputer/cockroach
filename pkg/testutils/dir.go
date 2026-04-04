// Copyright 2016 The Cockroach Authors.
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

package testutils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/build/bazel"
	"github.com/cockroachdb/cockroach/pkg/util/fileutil"
)

// TempDir creates a directory and a function to clean it up at the end of the
// test.
func TempDir(t testing.TB) (string, func()) {
	tmpDir := ""
	if bazel.BuiltWithBazel() {
		// Bazel sets up private temp directories for each test.
		// Normally, this private temp directory will be cleaned up automatically.
		// However, we do use external tools (such as stress) which re-execute the
		// same test multiple times.  Bazel, on the other hand, does not know about
		// this, and only creates this temporary directory once.  So, ensure we create
		// a unique temporary directory underneath bazel TEST_TMPDIR.
		tmpDir = bazel.TestTmpDir()
	}

	dir, err := ioutil.TempDir(tmpDir, fileutil.EscapeFilename(t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}
	return dir, cleanup
}
