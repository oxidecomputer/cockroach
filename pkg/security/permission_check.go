// Copyright 2017 The Cockroach Authors.
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

package security

import (
	"github.com/cockroachdb/cockroach/pkg/util/sysutil"
	"github.com/cockroachdb/errors"
)

// checkFilePermissions takes the passed path and file info, and returns an
// error if the file fails to match the permissions required. If this function
// returns nil, the file's permissions are acceptable.
func checkFilePermissions(processGID int, fullKeyPath string, fileACL sysutil.ACLInfo) error {
	// if the file is owned by root but also owned by the process's group
	// ID, we'll make an exception.
	if fileACL.IsOwnedByUID(uint64(0)) && fileACL.IsOwnedByGID(uint64(processGID)) {
		// if the file is owned by root, we allow those in the owning group to read it
		if sysutil.ExceedsPermissions(fileACL.Mode(), maxGroupKeyPermissions) {
			return errors.Errorf("key file %s has permissions %s, exceeds %s",
				fullKeyPath, fileACL.Mode(), maxGroupKeyPermissions)

		}

		return nil
	}

	if sysutil.ExceedsPermissions(fileACL.Mode(), maxKeyPermissions) {
		return errors.Errorf("key file %s has permissions %s, exceeds %s",
			fullKeyPath, fileACL.Mode(), maxKeyPermissions)
	}

	return nil
}
