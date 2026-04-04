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

package cli

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestFilePrefixerWithOptions(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	fileRE := regexp.MustCompile(logFilePattern)
	fInfo := func(path string, fileRE *regexp.Regexp) fileInfo {
		f, _ := getLogFileInfo(path, fileRE)
		return f
	}

	cases := []struct {
		name     string
		files    []fileInfo
		options  []filePrefixerOption
		expected [][]byte
	}{
		{
			"TestFilePrefixerWithOptions(default)",
			[]fileInfo{
				fInfo("nodes/1/cockroach.test-0001.ubuntu.2018-11-30T22_06_47Z.003959.log", fileRE),
				fInfo("nodes/2/cockroach.test-0001.ubuntu.2018-11-29T22_05_47Z.003955.log", fileRE),
				fInfo("nodes/9/cockroach.test-0001.ubuntu.2018-11-30T22_06_47Z.003960.log", fileRE),
			},
			[]filePrefixerOption{},
			[][]byte{

				[]byte("test-0001> "),
				[]byte("test-0001> "),
				[]byte("test-0001> "),
			},
		},
		{
			"TestFilePrefixerWithOptions(template)",
			[]fileInfo{
				fInfo("nodes/1/nodes/cockroach.test-0001.ubuntu.2018-11-30T22_06_47Z.003959.log", fileRE),
				fInfo("nodes/2/cockroach.test-0001.ubuntu.2018-11-29T22_05_47Z.003955.log", fileRE),
				fInfo("nodes/9/cockroach.test-0001.ubuntu.2018-11-30T22_06_47Z.003960.log", fileRE),
				fInfo("nodes/nodes/4/cockroach.test-0001.ubuntu.2018-11-30T22_06_47Z.003960.log", fileRE),
			},
			[]filePrefixerOption{withTemplate("(${fpath}) ${host}>")},
			[][]byte{

				[]byte("(1) test-0001>"),
				[]byte("(2) test-0001>"),
				[]byte("(9) test-0001>"),
				[]byte("(4) test-0001>"),
			},
		},
		{
			"TestFilePrefixerWithOptions(common tokens per file path)",
			[]fileInfo{
				fInfo("1/logs/1/cockroach.test-0001.ubuntu.2018-11-30T22_06_47Z.003959.log", fileRE),
				fInfo("2/nodes/2/cockroach.test-0001.ubuntu.2018-11-29T22_05_47Z.003955.log", fileRE),
			},
			[]filePrefixerOption{withTemplate("(${fpath}) ${host}>")},
			[][]byte{

				[]byte("(1/logs/1) test-0001>"),
				[]byte("(2/nodes/2) test-0001>"),
			},
		},
	}

	for _, tt := range cases {
		prefixer := newFilePrefixer(tt.options...)
		prefixer.PopulatePrefixes(tt.files)

		for i, file := range tt.files {
			if !bytes.Equal(file.prefix, tt.expected[i]) {
				t.Errorf("%s: got %v, expected %v for fp: %s",
					tt.name, string(file.prefix), string(tt.expected[i]), file.path)
			}
		}
	}
}
