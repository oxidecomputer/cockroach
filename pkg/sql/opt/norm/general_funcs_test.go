// Copyright 2019 The Cockroach Authors.
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

package norm

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

func TestCommuteJoinFlags(t *testing.T) {
	defer leaktest.AfterTest(t)()

	cases := [][2]memo.JoinFlags{
		{0, 0},

		{
			memo.DisallowLookupJoinIntoLeft,
			memo.DisallowLookupJoinIntoRight,
		},

		{
			memo.DisallowInvertedJoinIntoLeft,
			memo.DisallowInvertedJoinIntoRight,
		},

		{
			memo.PreferLookupJoinIntoLeft,
			memo.PreferLookupJoinIntoRight,
		},

		{
			memo.AllowOnlyMergeJoin,
			memo.AllowOnlyMergeJoin,
		},

		{
			memo.DisallowHashJoinStoreLeft | memo.DisallowMergeJoin | memo.DisallowLookupJoinIntoLeft | memo.DisallowLookupJoinIntoRight |
				memo.DisallowInvertedJoinIntoLeft | memo.DisallowInvertedJoinIntoRight,
			memo.DisallowHashJoinStoreRight | memo.DisallowMergeJoin | memo.DisallowLookupJoinIntoLeft | memo.DisallowLookupJoinIntoRight |
				memo.DisallowInvertedJoinIntoLeft | memo.DisallowInvertedJoinIntoRight,
		},

		{
			memo.DisallowHashJoinStoreLeft | memo.DisallowHashJoinStoreRight | memo.DisallowMergeJoin | memo.DisallowLookupJoinIntoLeft |
				memo.DisallowInvertedJoinIntoLeft | memo.DisallowInvertedJoinIntoRight,
			memo.DisallowHashJoinStoreLeft | memo.DisallowHashJoinStoreRight | memo.DisallowMergeJoin | memo.DisallowLookupJoinIntoRight |
				memo.DisallowInvertedJoinIntoLeft | memo.DisallowInvertedJoinIntoRight,
		},

		{
			memo.DisallowMergeJoin | memo.DisallowHashJoinStoreLeft | memo.DisallowLookupJoinIntoRight | memo.DisallowInvertedJoinIntoRight,
			memo.DisallowMergeJoin | memo.DisallowHashJoinStoreRight | memo.DisallowLookupJoinIntoLeft | memo.DisallowInvertedJoinIntoLeft,
		},
	}

	var funcs CustomFuncs
	for _, tc := range cases {
		// The result of commuting flags should be symmetrical, so test each case in
		// both directions.
		for dir := 0; dir <= 1; dir++ {
			in, out := tc[dir], tc[dir^1]
			res := funcs.CommuteJoinFlags(&memo.JoinPrivate{Flags: in})
			if res.Flags != out {
				t.Errorf("input: '%s'  expected: '%s'  got: '%s'", in, out, res.Flags)
			}
		}
	}
}
