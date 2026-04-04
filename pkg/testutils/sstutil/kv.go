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

package sstutil

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
)

// KV is a simplified representation of an MVCC key/value pair.
type KV struct {
	KeyString     string
	WallTimestamp int64  // 0 for inline
	ValueString   string // "" for nil (tombstone)
}

// Key returns the roachpb.Key representation of the key.
func (kv KV) Key() roachpb.Key {
	return roachpb.Key(kv.KeyString)
}

// Timestamp returns the hlc.Timestamp representation of the timestamp.
func (kv KV) Timestamp() hlc.Timestamp {
	return hlc.Timestamp{WallTime: kv.WallTimestamp}
}

// MVCCKey returns the storage.MVCCKey representation of the key and timestamp.
func (kv KV) MVCCKey() storage.MVCCKey {
	return storage.MVCCKey{
		Key:       kv.Key(),
		Timestamp: kv.Timestamp(),
	}
}

// Value returns the roachpb.Value representation of the value.
func (kv KV) Value() roachpb.Value {
	value := roachpb.MakeValueFromString(kv.ValueString)
	if kv.ValueString == "" {
		value = roachpb.Value{}
	}
	value.InitChecksum(kv.Key())
	return value
}

// ValueBytes returns the roachpb.Value byte-representation of the value.
func (kv KV) ValueBytes() []byte {
	return kv.Value().RawBytes
}
