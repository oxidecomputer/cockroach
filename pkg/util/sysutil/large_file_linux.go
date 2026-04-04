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

//go:build linux
// +build linux

package sysutil

import (
	"os"
	"syscall"

	"github.com/cockroachdb/errors"
	"golang.org/x/sys/unix"
)

// ResizeLargeFile resizes the file at the given path to be the provided
// length in bytes. If no file exists at path, ResizeLargeFile creates a file.
// All disk blocks within the new file are allocated, and there are no sparse
// regions.
//
// On Linux, it uses the fallocate syscall to efficiently allocate disk space.
// On other platforms, it naively writes the specified number of bytes, which
// can take a long time.
func ResizeLargeFile(path string, bytes int64) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// Ensure that at least length bytes of the file are allocated. If the
	// file already existed, the disk blocks may not have been allocated even
	// if the file size is greater than length.
	if bytes > 0 {
		err := unix.Fallocate(int(f.Fd()), 0, 0, bytes)

		// Some linux filesystems, like older versions of ZFS, do not
		// support fallocate. If an error indicates it's not supported,
		// fallback to the naive implementation.
		var errno syscall.Errno
		if err != nil && errors.As(err, &errno) && (errno == syscall.ENOTSUP || errno == syscall.EOPNOTSUPP) {
			return resizeLargeFileNaive(path, bytes)
		} else if err != nil {
			return errors.Wrap(err, "fallocate")
		}
	}

	// Truncate down to bytes, in case the file is longer than bytes. This
	// will be a no-op if the file is already at the desired length.
	if err := unix.Ftruncate(int(f.Fd()), bytes); err != nil {
		return errors.Wrap(err, "ftruncate")
	}
	return errors.Wrap(f.Sync(), "fsync")
}
