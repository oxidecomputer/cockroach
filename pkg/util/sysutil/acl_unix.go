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

//go:build !windows && !plan9
// +build !windows,!plan9

package sysutil

import (
	"os"
	"syscall"
)

type unixACLInfo struct {
	uid  uint64
	gid  uint64
	mode os.FileMode
}

func (acl *unixACLInfo) UID() uint64 {
	return acl.uid
}

func (acl *unixACLInfo) GID() uint64 {
	return acl.gid
}

func (acl *unixACLInfo) IsOwnedByUID(uid uint64) bool {
	return acl.uid == uid
}

func (acl *unixACLInfo) IsOwnedByGID(gid uint64) bool {
	return acl.gid == gid
}

func (acl *unixACLInfo) Mode() os.FileMode {
	return acl.mode
}

// GetFileACLInfo returns an ACLInfo that has the UID and GID populated from
// the system specific file information.
func GetFileACLInfo(info os.FileInfo) ACLInfo {
	sysInfo := info.Sys()
	if nil != sysInfo {
		if statTInfo, ok := sysInfo.(*syscall.Stat_t); ok {
			// we use uint64 because a process should never have a
			// GID that's less than zero, and syscall.Stat_t may
			// use a uint16 or uint32 for GID, since it's often a
			// C struct under the hood.
			return &unixACLInfo{
				uid:  uint64(statTInfo.Uid),
				gid:  uint64(statTInfo.Gid),
				mode: info.Mode().Perm(),
			}
		}
	}

	// if we don't know who owns a file, assume that root owns it.
	return &unixACLInfo{
		uid:  uint64(0),
		gid:  uint64(0),
		mode: info.Mode().Perm(),
	}
}
