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

package tree_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/faketreeeval"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/randgen"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/lib/pq/oid"
)

// TestCastMap tests that every cast in tree.castMap can be performed by
// PerformCast.
func TestCastMap(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	evalCtx := tree.MakeTestingEvalContext(cluster.MakeTestingClusterSettings())
	rng, _ := randutil.NewTestRand()
	evalCtx.Planner = &faketreeeval.DummyEvalPlanner{}

	tree.ForEachCast(func(src, tgt oid.Oid) {
		srcType := types.OidToType[src]
		tgtType := types.OidToType[tgt]
		srcDatum := randgen.RandDatum(rng, srcType, false /* nullOk */)

		// TODO(mgartner): We do not allow casting a negative integer to bit
		// types with unbounded widths. Until we add support for this, we
		// ensure that the srcDatum is positive.
		if srcType.Family() == types.IntFamily && tgtType.Family() == types.BitFamily {
			srcVal := *srcDatum.(*tree.DInt)
			if srcVal < 0 {
				srcDatum = tree.NewDInt(-srcVal)
			}
		}

		_, err := tree.PerformCast(&evalCtx, srcDatum, tgtType)
		// If the error is a CannotCoerce error, then PerformCast does not
		// support casting from src to tgt. The one exception is negative
		// integers to bit types which return the same error code (see the TODO
		// above).
		if err != nil && pgerror.HasCandidateCode(err) && pgerror.GetPGCode(err) == pgcode.CannotCoerce {
			t.Errorf("cast from %s to %s failed: %s", srcType, tgtType, err)
		}
	})
}
