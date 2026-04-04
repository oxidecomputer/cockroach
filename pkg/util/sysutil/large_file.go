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
	"os"

	"github.com/cockroachdb/errors"
)

func resizeLargeFileNaive(path string, bytes int64) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	sixtyFourMB := make([]byte, 64<<20)
	for bytes > 0 {
		z := sixtyFourMB
		if bytes < int64(len(z)) {
			z = sixtyFourMB[:bytes]
		}
		if _, err := f.Write(z); err != nil {
			return errors.Wrap(err, "write")
		}
		bytes -= int64(len(z))
	}
	return errors.Wrap(f.Sync(), "fsync")
}
