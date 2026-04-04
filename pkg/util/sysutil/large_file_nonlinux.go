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

//go:build !linux
// +build !linux

package sysutil

// ResizeLargeFile resizes the file at the given path to be the provided
// length in bytes. If no file exists at path, ResizeLargeFile creates a file.
// All disk blocks within the new file are allocated, and there are no sparse
// regions.
//
// On Linux, it uses the fallocate syscall to efficiently allocate disk space.
// On other platforms, it naively writes the specified number of bytes, which
// can take a long time.
func ResizeLargeFile(path string, bytes int64) error {
	return resizeLargeFileNaive(path, bytes)
}
