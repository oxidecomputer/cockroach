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

package colcontainer_test

import (
	"context"
	"os"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/col/coldataext"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/cockroach/pkg/util/mon"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
)

var (
	// testAllocator is a colmem.Allocator with an unlimited budget for use in
	// tests.
	testAllocator     *colmem.Allocator
	testColumnFactory coldata.ColumnFactory

	// testMemMonitor and testMemAcc are a test monitor with an unlimited budget
	// and a memory account bound to it for use in tests.
	testMemMonitor *mon.BytesMonitor
	testMemAcc     *mon.BoundAccount

	// testDiskMonitor and testDiskmAcc are a test monitor with an unlimited budget
	// and a disk account bound to it for use in tests.
	testDiskMonitor *mon.BytesMonitor
	testDiskAcc     *mon.BoundAccount
)

func TestMain(m *testing.M) {
	randutil.SeedForTests()
	os.Exit(func() int {
		ctx := context.Background()
		testMemMonitor = execinfra.NewTestMemMonitor(ctx, cluster.MakeTestingClusterSettings())
		defer testMemMonitor.Stop(ctx)
		memAcc := testMemMonitor.MakeBoundAccount()
		testMemAcc = &memAcc
		testColumnFactory = coldataext.NewExtendedColumnFactory(nil /* evalCtx */)
		testAllocator = colmem.NewAllocator(ctx, testMemAcc, testColumnFactory)
		defer testMemAcc.Close(ctx)

		testDiskMonitor = execinfra.NewTestDiskMonitor(ctx, cluster.MakeTestingClusterSettings())
		defer testDiskMonitor.Stop(ctx)
		diskAcc := testDiskMonitor.MakeBoundAccount()
		testDiskAcc = &diskAcc
		defer testDiskAcc.Close(ctx)
		return m.Run()
	}())
}
