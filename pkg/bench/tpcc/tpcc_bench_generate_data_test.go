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

package tpcc

import (
	"flag"
	"os"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/stretchr/testify/require"
)

var (
	storeDirFlag = flag.String(
		"store-dir", "", "name of the directory on disk to use for the loaded TPCC state",
	)
	generateStoreDirFlag = flag.Bool("generate-store-dir", false,
		"if store-dir is set, remove any exist data and regenerate the data")
)

func maybeGenerateStoreDir(b testing.TB) (_ vfs.FS, storeDir string, cleanup func()) {
	storeDir = *storeDirFlag
	cleanup = func() {}

	if storeDir != "" {
		if !*generateStoreDirFlag {
			fi, err := os.Stat(storeDir)
			require.NoError(b, err, "consider --generate-store-dir")
			require.True(b, fi.IsDir(), "consider --generate-store-dir")
			return vfs.Default, *storeDirFlag, cleanup
		}
		require.NoError(b, os.RemoveAll(storeDir))
		require.NoError(b, os.MkdirAll(storeDir, 0777))
	} else {
		storeDir, cleanup = testutils.TempDir(b)
		defer func() {
			if b.Failed() {
				cleanup()
			}
		}()
	}

	require.NoError(b, generateStoreDir.exec(cmdEnv{
		{storeDirEnvVar, storeDir},
	}).Run())
	return vfs.Default, storeDir, cleanup
}
