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

package rangefeedbuffer

import (
	"context"
	"math"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/stretchr/testify/require"
)

// TestMergeKVs tests the logic of MergeKVs and the logic of to extract
// KVs from a slice of events.
func TestMergeKVs(t *testing.T) {
	type row struct {
		key   string
		ts    int64
		value string
	}
	prefix := keys.SystemSQLCodec.TablePrefix(1)
	prefix = prefix[:len(prefix):len(prefix)]
	mkKey := func(r row) roachpb.Key {
		return encoding.EncodeStringAscending(prefix, r.key)
	}
	toKeyValue := func(r row) (kv roachpb.KeyValue) {
		kv.Key = mkKey(r)
		kv.Value.Timestamp = hlc.Timestamp{WallTime: r.ts}
		if r.value != "" {
			kv.Value.SetString(r.value)
		}
		kv.Value.InitChecksum(kv.Key)
		return kv
	}
	toRangeFeedEvent := func(r row) *roachpb.RangeFeedValue {
		kv := toKeyValue(r)
		return &roachpb.RangeFeedValue{
			Key:   kv.Key,
			Value: kv.Value,
		}
	}
	toKVs := func(rows []row) (kvs []roachpb.KeyValue) {
		for _, r := range rows {
			kvs = append(kvs, toKeyValue(r))
		}
		return kvs
	}
	toBuffer := func(t *testing.T, rows []row) *Buffer {
		buf := New(len(rows))
		for _, r := range rows {
			require.NoError(t, buf.Add(toRangeFeedEvent(r)))
		}
		return buf
	}
	toKVsThroughBuffer := func(t *testing.T, rows []row) []roachpb.KeyValue {
		return EventsToKVs(toBuffer(t, rows).Flush(
			context.Background(),
			hlc.Timestamp{WallTime: math.MaxInt64},
		), RangeFeedValueEventToKV)
	}
	type testCase [3][]row // (a, b, merged)
	for _, tc := range []testCase{
		{
			{
				{"a", 1, "asdf"},
				{"a", 4, "af"},
				{"c", 3, "boo"},
			},
			{
				{"a", 5, "winner"},
				{"b", 5, "2as"},
				{"b", 1, "2"},
				{"c", 4, ""},
				{"d", 4, ""},
			},
			{
				{"a", 5, "winner"},
				{"b", 5, "2as"},
			},
		},
	} {
		a, b, merged := tc[0], tc[1], tc[2]

		require.Equal(t, toKVs(merged), MergeKVs(toKVs(a), toKVs(b)))

		// Exercise the conversions out of RangeFeedValue.
		require.Equal(t, toKVs(merged), MergeKVs(
			toKVsThroughBuffer(t, a), toKVsThroughBuffer(t, b),
		))
	}
}
