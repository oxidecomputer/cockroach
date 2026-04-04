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

package tpch

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/bufalloc"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestRandPartName(t *testing.T) {
	var a bufalloc.ByteAllocator
	rng := rand.New(rand.NewSource(uint64(timeutil.Now().UnixNano())))
	seen := make(map[string]int)
	runOneRound := func() {
		res := randPartName(rng, &a)
		names := strings.Split(string(res), " ")
		assert.Equal(t, len(names), nPartNames)
		seenLocal := make(map[string]int)
		for _, name := range names {
			if _, ok := seenLocal[name]; ok {
				t.Errorf("names in '%s' are not unique", res)
			}
			seenLocal[name]++
			seen[name]++
		}
	}

	// We can't guarantee much about the global distribution of names,
	// but we should make sure that we're not always using the same 5
	// names. Run up to 100 times before failing.
	//
	// NB: The odds of this flaking are extremely low. 92 choose 5 gives
	// 4,9177,128 unique combinations. After 100 shuffles, the probability of
	// seeing the same combination is astronomically low.
	for i := 0; i < 100; i++ {
		if len(seen) > nPartNames {
			return
		}
		runOneRound()
	}

	if len(seen) <= nPartNames {
		t.Errorf("only saw %d names after calling randPartName 100 times", nPartNames)
	}
}
