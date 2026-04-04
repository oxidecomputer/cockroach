// Copyright 2018 The Cockroach Authors.
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

package sysutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestResizeLargeFile(t *testing.T) {
	d, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(d); err != nil {
			t.Fatal(err)
		}
	}()
	fname := filepath.Join(d, "ballast")

	lens := []int64{2000, 1000, 64<<20 + 10, 0, 1}
	for _, n := range lens {
		if err := ResizeLargeFile(fname, n); err != nil {
			t.Fatal(err)
		}
		fi, err := os.Stat(fname)
		if err != nil {
			t.Fatal(err)
		}
		if n != fi.Size() {
			t.Fatalf("expected size of file %d, got %d", n, fi.Size())
		}
	}
}
