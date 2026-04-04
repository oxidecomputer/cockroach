// Copyright 2015 The Cockroach Authors.
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

package catalogkeys

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

func TestKeyAddress(t *testing.T) {
	defer leaktest.AfterTest(t)()
	tenSysCodec := keys.SystemSQLCodec
	ten5Codec := keys.MakeSQLCodec(roachpb.MakeTenantID(5))
	testCases := []struct {
		key roachpb.Key
	}{
		{MakeDescMetadataKey(tenSysCodec, 123)},
		{MakeDescMetadataKey(tenSysCodec, 124)},
		{MakePublicObjectNameKey(tenSysCodec, 0, "BAR")},
		{MakePublicObjectNameKey(tenSysCodec, 1, "BAR")},
		{MakePublicObjectNameKey(tenSysCodec, 1, "foo")},
		{MakePublicObjectNameKey(tenSysCodec, 2, "foo")},
		{MakeDescMetadataKey(ten5Codec, 123)},
		{MakeDescMetadataKey(ten5Codec, 124)},
		{MakePublicObjectNameKey(ten5Codec, 0, "BAR")},
		{MakePublicObjectNameKey(ten5Codec, 1, "BAR")},
		{MakePublicObjectNameKey(ten5Codec, 1, "foo")},
		{MakePublicObjectNameKey(ten5Codec, 2, "foo")},
	}
	var lastKey roachpb.Key
	for i, test := range testCases {
		resultAddr, err := keys.Addr(test.key)
		if err != nil {
			t.Fatal(err)
		}
		result := resultAddr.AsRawKey()
		if result.Compare(lastKey) <= 0 {
			t.Errorf("%d: key address %q is <= %q", i, result, lastKey)
		}
		lastKey = result
	}
}
