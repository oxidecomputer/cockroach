// Copyright 2018 The Cockroach Authors.
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

package rowexec

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/randgen"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/testutils/distsqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

// TestInputStatCollector verifies that an inputStatCollector correctly collects
// stats from an input.
func TestInputStatCollector(t *testing.T) {
	defer leaktest.AfterTest(t)()

	const numRows = 100

	isc := newInputStatCollector(
		distsqlutils.NewRowBuffer(types.OneIntCol, randgen.MakeIntRows(numRows, 1), distsqlutils.RowBufferArgs{}),
	)
	for row, meta := isc.Next(); row != nil || meta != nil; row, meta = isc.Next() {
	}
	if isc.stats.NumTuples.Value() != numRows {
		t.Fatalf("counted %s rows but expected %d", isc.stats.NumTuples, numRows)
	}
}
