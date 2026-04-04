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
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/pebble/vfs"
)

func runTestSSTIterator(t *testing.T, iter SimpleMVCCIterator, allKVs []MVCCKeyValue) {
	// Drop the first kv so we can test Seek.
	expected := allKVs[1:]

	// Run the test multiple times to check re-Seeking.
	for i := 0; i < 3; i++ {
		var kvs []MVCCKeyValue
		for iter.SeekGE(expected[0].Key); ; iter.Next() {
			ok, err := iter.Valid()
			if err != nil {
				t.Fatalf("%+v", err)
			}
			if !ok {
				break
			}
			kv := MVCCKeyValue{
				Key: MVCCKey{
					Key:       append([]byte(nil), iter.UnsafeKey().Key...),
					Timestamp: iter.UnsafeKey().Timestamp,
				},
				Value: append([]byte(nil), iter.UnsafeValue()...),
			}
			kvs = append(kvs, kv)
		}
		if ok, err := iter.Valid(); err != nil {
			t.Fatalf("%+v", err)
		} else if ok {
			t.Fatalf("expected !ok")
		}

		lastElemKey := expected[len(expected)-1].Key
		seekTo := MVCCKey{Key: lastElemKey.Key.Next()}

		iter.SeekGE(seekTo)
		if ok, err := iter.Valid(); err != nil {
			t.Fatalf("%+v", err)
		} else if ok {
			foundKey := iter.UnsafeKey()
			t.Fatalf("expected !ok seeking to lastEmem.Next(). foundKey %s < seekTo %s: %t",
				foundKey, seekTo, foundKey.Less(seekTo))
		}

		if !reflect.DeepEqual(kvs, expected) {
			t.Fatalf("got %+v but expected %+v", kvs, expected)
		}
	}
}

func TestSSTIterator(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	st := cluster.MakeTestingClusterSettings()
	sstFile := &MemFile{}
	sst := MakeIngestionSSTWriter(ctx, st, sstFile)
	defer sst.Close()
	var allKVs []MVCCKeyValue
	for i := 0; i < 10; i++ {
		kv := MVCCKeyValue{
			Key: MVCCKey{
				Key:       []byte{'A' + byte(i)},
				Timestamp: hlc.Timestamp{WallTime: int64(i)},
			},
			Value: []byte{'a' + byte(i)},
		}
		if err := sst.Put(kv.Key, kv.Value); err != nil {
			t.Fatalf("%+v", err)
		}
		allKVs = append(allKVs, kv)
	}

	if err := sst.Finish(); err != nil {
		t.Fatalf("%+v", err)
	}

	t.Run("Disk", func(t *testing.T) {
		tempDir, cleanup := testutils.TempDir(t)
		defer cleanup()

		path := filepath.Join(tempDir, "data.sst")
		if err := ioutil.WriteFile(path, sstFile.Data(), 0600); err != nil {
			t.Fatalf("%+v", err)
		}
		file, err := vfs.Default.Open(path)
		if err != nil {
			t.Fatal(err)
		}

		iter, err := NewSSTIterator(file)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer iter.Close()
		runTestSSTIterator(t, iter, allKVs)
	})
	t.Run("Mem", func(t *testing.T) {
		iter, err := NewMemSSTIterator(sstFile.Data(), false)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer iter.Close()
		runTestSSTIterator(t, iter, allKVs)
	})
}
