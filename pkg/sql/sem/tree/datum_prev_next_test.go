// Copyright 2023 The Cockroach Authors.
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

package tree_test

import (
	"math"
	"testing"

	"github.com/cockroachdb/apd/v3"
	"github.com/cockroachdb/cockroach/pkg/sql/randgen"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/stretchr/testify/require"
)

// TestDatumPrevNext verifies that tree.DatumPrev and tree.DatumNext return
// datums that are smaller and larger, respectively, than the given datum if
// ok=true is returned (modulo some edge cases).
func TestDatumPrevNext(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	rng, _ := randutil.NewTestRand()
	var evalCtx tree.EvalContext
	const numRuns = 1000
	for i := 0; i < numRuns; i++ {
		typ := randgen.RandType(rng)
		d := randgen.RandDatum(rng, typ, false /* nullOk */)
		// Ignore NaNs and infinities.
		if f, ok := d.(*tree.DFloat); ok {
			if math.IsNaN(float64(*f)) || math.IsInf(float64(*f), 0) {
				continue
			}
		}
		if dec, ok := d.(*tree.DDecimal); ok {
			if dec.Form == apd.NaN || dec.Form == apd.Infinite {
				continue
			}
		}
		if !d.IsMin(&evalCtx) {
			if prev, ok := tree.DatumPrev(d, &evalCtx, &evalCtx.CollationEnv); ok {
				cmp, err := d.CompareError(&evalCtx, prev)
				require.NoError(t, err)
				require.True(t, cmp > 0, "d=%s, prev=%s, type=%s", d.String(), prev.String(), d.ResolvedType().SQLString())
			}
		}
		if !d.IsMax(&evalCtx) {
			if next, ok := tree.DatumNext(d, &evalCtx, &evalCtx.CollationEnv); ok {
				cmp, err := d.CompareError(&evalCtx, next)
				require.NoError(t, err)
				require.True(t, cmp < 0, "d=%s, next=%s, type=%s", d.String(), next.String(), d.ResolvedType().SQLString())
			}
		}
	}
}
