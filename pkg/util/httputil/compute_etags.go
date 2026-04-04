// Copyright 2022 The Cockroach Authors.
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

package httputil

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/fs"

	"github.com/cockroachdb/errors"
)

// ComputeEtags recursively computes the SHA1 hash of every file in fsys
// starting at root, and stores the computed hashes in dest["/path/to/file"]
// (*including* a leading "/") for later use as an ETag response header.
func ComputeEtags(fsys fs.FS, dest map[string]string) error {
	if dest == nil {
		return errors.New("Unable to hash files without a hash destination")
	}

	hash := sha1.New()

	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		hash.Reset()

		file, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		// Copy file contents into the hash algorithm
		if _, err := io.Copy(hash, file); err != nil {
			return err
		}

		// Store the computed hash
		dest["/"+path] = hex.EncodeToString(hash.Sum(nil))
		return nil
	})
}
