// Copyright 2020 The Cockroach Authors.
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

package kvnemesis

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	k := func(s string, ts hlc.Timestamp) storage.MVCCKey {
		return storage.MVCCKey{Key: []byte(s), Timestamp: ts}
	}
	var missing roachpb.Value
	v := func(s string, ts hlc.Timestamp) roachpb.Value {
		v := roachpb.MakeValueFromString(s)
		v.Timestamp = ts
		return v
	}
	ts := func(i int) hlc.Timestamp {
		return hlc.Timestamp{WallTime: int64(i)}
	}

	e, err := MakeEngine()
	require.NoError(t, err)
	defer e.Close()
	assert.Equal(t, missing, e.Get(roachpb.Key(`a`), ts(1)))
	e.Put(k(`a`, ts(1)), roachpb.MakeValueFromString(`a-1`).RawBytes)
	e.Put(k(`a`, ts(2)), roachpb.MakeValueFromString(`a-2`).RawBytes)
	e.Put(k(`b`, ts(2)), roachpb.MakeValueFromString(`b-2`).RawBytes)
	e.Delete(k(`b`, ts(3)))
	e.Delete(k(`c`, ts(4)))
	e.Put(k(`d`, ts(4)), roachpb.MakeValueFromString(`d-4`).RawBytes)
	e.Put(k(`e`, ts(4)), roachpb.MakeValueFromString(`e-4`).RawBytes)
	e.Delete(k(`d`, ts(5)))
	assert.Equal(t, v(`a-2`, ts(2)), e.Get(roachpb.Key(`a`), ts(3)))
	assert.Equal(t, v(`a-2`, ts(2)), e.Get(roachpb.Key(`a`), ts(2)))
	assert.Equal(t, v(`a-1`, ts(1)), e.Get(roachpb.Key(`a`), ts(2).Prev()))
	assert.Equal(t, v(`a-1`, ts(1)), e.Get(roachpb.Key(`a`), ts(1)))
	assert.Equal(t, missing, e.Get(roachpb.Key(`a`), ts(1).Prev()))
	assert.Equal(t, v(`b-2`, ts(2)), e.Get(roachpb.Key(`b`), ts(3).Prev()))
	assert.Equal(t, v(`b-2`, ts(2)), e.Get(roachpb.Key(`b`), ts(2)))
	assert.Equal(t, missing, e.Get(roachpb.Key(`b`), ts(1)))
	assert.Equal(t, missing, e.Get(roachpb.Key(`b`), ts(4)))
	assert.Equal(t, v(`d-4`, ts(4)), e.Get(roachpb.Key(`d`), ts(4)))
	assert.Equal(t, missing, e.Get(roachpb.Key(`d`), ts(5)))
	assert.Equal(t, v(`e-4`, ts(4)), e.Get(roachpb.Key(`e`), ts(5)))

	assert.Equal(t, strings.TrimSpace(`
"a" 0.000000002,0 -> /BYTES/a-2
"a" 0.000000001,0 -> /BYTES/a-1
"b" 0.000000003,0 -> /<empty>
"b" 0.000000002,0 -> /BYTES/b-2
"c" 0.000000004,0 -> /<empty>
"d" 0.000000005,0 -> /<empty>
"d" 0.000000004,0 -> /BYTES/d-4
"e" 0.000000004,0 -> /BYTES/e-4
	`), e.DebugPrint(""))
}
