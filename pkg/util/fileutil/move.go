// Copyright 2019 The Cockroach Authors.
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

package fileutil

import (
	"os"

	"github.com/cockroachdb/cockroach/pkg/util/sysutil"
	"github.com/cockroachdb/errors"
)

// Move moves a file from a directory to another, while handling
// cross-filesystem moves properly.
// If the target file already exists, it is truncated.
// If the move fails, then the target file may be left in an inconsistent state.
func Move(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if !isCrossDeviceLinkError(err) {
		return err
	}

	if err = CopyFile(oldPath, newPath); err != nil {
		return err
	}

	return os.RemoveAll(oldPath)
}

func isCrossDeviceLinkError(err error) bool {
	if err == nil {
		return false
	}
	var le *os.LinkError
	if errors.As(err, &le) {
		return sysutil.IsCrossDeviceLinkErrno(le.Err)
	}
	return false
}
