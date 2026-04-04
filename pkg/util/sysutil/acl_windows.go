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

type windowsACLInfo struct {
	mode os.FileMode
}

func (acl *windowsACLInfo) UID() uint64 {
	return uint64(0)
}

func (acl *windowsACLInfo) GID() uint64 {
	return uint64(0)
}

func (acl *windowsACLInfo) IsOwnedByUID(uid uint64) bool {
	return acl.UID() == uid
}

func (acl *windowsACLInfo) IsOwnedByGID(gid uint64) bool {
	return acl.GID() == gid
}

func (acl *windowsACLInfo) Mode() os.FileMode {
	return acl.mode
}

// GetFileACLInfo returns an ACLInfo that has the UID and GID populated from
// the system specific file information. On Windows, this returns an ACLInfo
// that always has the UID and GID of 0, since we haven't implemented support
// for looking up the Windows owner information.
func GetFileACLInfo(info os.FileInfo) ACLInfo {
	return &windowsACLInfo{
		mode: info.Mode().Perm(),
	}
}
