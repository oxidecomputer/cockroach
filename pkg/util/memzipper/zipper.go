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

package memzipper

import (
	"archive/zip"
	"bytes"

	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

// Zipper builds a zip file into an in-memory buffer.
type Zipper struct {
	buf *bytes.Buffer
	z   *zip.Writer
	err error
}

// Init initializes the underlying Zipper with a new zip writer.
func (z *Zipper) Init() {
	z.buf = &bytes.Buffer{}
	z.z = zip.NewWriter(z.buf)
}

// AddFile adds a file to the underlying Zipper.
func (z *Zipper) AddFile(name string, contents string) {
	if z.err != nil {
		return
	}
	w, err := z.z.CreateHeader(&zip.FileHeader{
		Name:     name,
		Method:   zip.Deflate,
		Modified: timeutil.Now(),
	})
	if err != nil {
		z.err = err
		return
	}
	_, z.err = w.Write([]byte(contents))
}

// Finalize finalizes the Zipper by closing the zip writer.
func (z *Zipper) Finalize() (*bytes.Buffer, error) {
	if z.err != nil {
		return nil, z.err
	}
	if err := z.z.Close(); err != nil {
		return nil, err
	}
	buf := z.buf
	*z = Zipper{}
	return buf, nil
}
