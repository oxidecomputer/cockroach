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

package scdecomp_test

import (
	"context"
	gosql "database/sql"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scrun"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/sctest"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/testcluster"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestDecomposeToElements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	ctx := context.Background()

	newCluster := func(t *testing.T, knobs *scrun.TestingKnobs) (_ *gosql.DB, cleanup func()) {
		tc := testcluster.StartTestCluster(t, 1, base.TestClusterArgs{})
		return tc.ServerConn(0), func() { tc.Stopper().Stop(ctx) }
	}

	sctest.DecomposeToElements(t, testutils.TestDataPath(t), newCluster)
}
