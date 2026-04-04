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

// Given the path to a directory and an output path, this executable creates a
// .tar.gz archive. This is not feature-complete (at all) compared to the `tar`
// utility, but works for the purposes we have (i.e. packaging UI assets).

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/cli/exit"
	"github.com/cockroachdb/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		exit.WithCode(exit.UnspecifiedError())
	}
}

func run() error {
	if len(os.Args) != 3 {
		return errors.Newf("usage: %s SRCDIR OUTFILE\n", os.Args[0])
	}
	os.Args[1] = strings.TrimRight(os.Args[1], "/")

	// Make tar archive
	var tarContents bytes.Buffer
	tarWriter := tar.NewWriter(&tarContents)
	err := filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == os.Args[1] && d.IsDir() {
			return nil
		}
		if d.IsDir() {
			return errors.Newf("cannot compress subdirectory %s", path)
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		info, err := src.Stat()
		if err != nil {
			return err
		}
		err = tarWriter.WriteHeader(&tar.Header{Name: d.Name(), Size: info.Size()})
		if err != nil {
			return errors.Wrap(err, "could not write header to tar file")
		}
		_, err = io.Copy(tarWriter, src)
		if err != nil {
			return err
		}
		return src.Close()
	})
	if err != nil {
		return err
	}
	err = tarWriter.Close()
	if err != nil {
		return err
	}

	// compress tar archive w/ gzip
	outFile, err := os.Create(os.Args[2])
	if err != nil {
		return err
	}
	gzipWriter := gzip.NewWriter(outFile)
	_, err = gzipWriter.Write(tarContents.Bytes())
	if err != nil {
		return err
	}
	err = gzipWriter.Close()
	if err != nil {
		return err
	}
	return outFile.Close()
}
