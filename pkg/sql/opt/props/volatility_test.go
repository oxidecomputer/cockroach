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

package props

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/stretchr/testify/require"
)

func TestVolatilitySet(t *testing.T) {
	var v VolatilitySet

	check := func(str string, isLeakProof, hasStable, hasVolatile bool) {
		t.Helper()

		require.Equal(t, v.String(), str)
		require.Equal(t, v.IsLeakProof(), isLeakProof)
		require.Equal(t, v.HasStable(), hasStable)
		require.Equal(t, v.HasVolatile(), hasVolatile)
	}
	check("leak-proof", true, false, false)

	v.Add(tree.VolatilityLeakProof)
	check("leak-proof", true, false, false)

	v.AddImmutable()
	check("immutable", false, false, false)

	v.AddStable()
	check("stable", false, true, false)

	v.AddVolatile()
	check("stable+volatile", false, true, true)

	v = 0
	v.AddVolatile()
	check("volatile", false, false, true)

	var w VolatilitySet
	w.AddImmutable()
	v.UnionWith(w)
	check("volatile", false, false, true)

	w.AddStable()
	v.UnionWith(w)
	check("stable+volatile", false, true, true)
}
