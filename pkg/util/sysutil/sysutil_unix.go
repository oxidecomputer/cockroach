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

//go:build !windows
// +build !windows

//lint:file-ignore Unconvert (redundant conversions are necessary for cross-platform compatibility)

package sysutil

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

// ProcessIdentity returns a string describing the user and group that this
// process is running as.
func ProcessIdentity() string {
	return fmt.Sprintf("uid %d euid %d gid %d egid %d",
		unix.Getuid(), unix.Geteuid(), unix.Getgid(), unix.Getegid())
}

// IsCrossDeviceLinkErrno checks whether the given error object (as
// extracted from an *os.LinkError) is a cross-device link/rename
// error.
func IsCrossDeviceLinkErrno(errno error) bool {
	return errno == syscall.EXDEV
}
