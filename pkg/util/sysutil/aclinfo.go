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

package sysutil

import "os"

// ACLInfo represents access control information for a file in the
// filesystem. It can be used to determine if a file has the correct
// permissions associated with it.
type ACLInfo interface {
	// UID returns the ID of the user that owns the file in question.
	UID() uint64
	// GID returns the ID of the group that owns the file in question.
	GID() uint64

	// IsOwnedByUID returns true of the passed UID is equal to the UID of the file in question.
	IsOwnedByUID(uint64) bool
	// IsOwnedByGID returns true of the passed GID is equal to the GID of the file in question.
	IsOwnedByGID(uint64) bool

	// Mode returns the os.FileMode representing the files permissions. Implementers of this
	// interface should return only the result of the `Perm` function in os.FileMode, to
	// ensure that only permissions are returned.
	Mode() os.FileMode
}

// ExceedsPermissions returns true if the passed os.FileMode represents a more stringent
// set of permissions than the ACLInfo's mode.
//
// For example, calling this function with an objectMode of 0640 will return false if
// the passed allowedMode is 0600.
func ExceedsPermissions(objectMode, allowedMode os.FileMode) bool {
	mask := os.FileMode(0777) ^ allowedMode
	return mask&objectMode != 0
}
