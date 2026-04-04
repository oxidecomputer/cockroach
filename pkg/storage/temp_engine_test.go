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

package storage

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestNewPebbleTempEngine(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	// Create the temporary directory for the RocksDB engine.
	tempDir, tempDirCleanup := testutils.TempDir(t)
	defer tempDirCleanup()

	db, fs, err := NewPebbleTempEngine(context.Background(), base.TempStorageConfig{Path: tempDir}, base.StoreSpec{Path: tempDir})
	if err != nil {
		t.Fatalf("error encountered when invoking NewRocksDBTempEngine: %+v", err)
	}
	defer db.Close()

	// Temp engine initialized with the temporary directory.
	if dir := fs.(*Pebble).path; tempDir != dir {
		t.Fatalf("temp engine initialized with unexpected parent directory.\nexpected %s\nactual %s",
			tempDir, dir)
	}
}
