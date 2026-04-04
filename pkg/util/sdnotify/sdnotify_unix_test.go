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

//go:build !windows
// +build !windows

package sdnotify

import (
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/cockroachdb/cockroach/pkg/util/log" // for flags
	"github.com/stretchr/testify/require"
)

func TestSDNotify(t *testing.T) {
	tmpDir := os.TempDir()
	// On BSD, binding to a socket is limited to a path length of 104 characters
	// (including the NUL terminator). In glibc, this limit is 108 characters.
	// macOS also has a tendency to produce very long temporary directory names.
	if len(tmpDir) >= 104-1-len("sdnotify/notify.sock")-10 {
		// Perhaps running inside a sandbox?
		t.Logf("default temp dir name is too long: %s", tmpDir)
		t.Logf("forcing use of /tmp instead")
		// Note: Using /tmp may fail on some systems; this is why we
		// prefer os.TempDir() by default.
		tmpDir = "/tmp"
	}
	dir, err := ioutil.TempDir(tmpDir, "sdnotify")
	require.NoError(t, err)

	l, err := listen(dir)
	require.NoError(t, err)
	defer func() { _ = l.close() }()

	ch := make(chan error)
	go func() {
		ch <- l.wait()
	}()

	if err := notify(l.Path, readyMsg); err != nil {
		t.Fatal(err)
	}
	if err := <-ch; err != nil {
		t.Fatal(err)
	}
}
