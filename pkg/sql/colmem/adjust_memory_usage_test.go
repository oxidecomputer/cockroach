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

package colmem

import (
	"context"
	"math"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/col/coldataext"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/sql/execinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/mon"
	"github.com/cockroachdb/redact"
	"github.com/stretchr/testify/require"
)

func TestAdjustMemoryUsage(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	st := cluster.MakeTestingClusterSettings()
	unlimitedMemMonitor := execinfra.NewTestMemMonitor(ctx, st)
	defer unlimitedMemMonitor.Stop(ctx)
	unlimitedMemAcc := unlimitedMemMonitor.MakeBoundAccount()
	defer unlimitedMemAcc.Close(ctx)

	limitedMemMonitorName := "test-limited"
	limit := int64(100000)
	limitedMemMonitor := mon.NewMonitorInheritWithLimit(
		redact.RedactableString(limitedMemMonitorName), limit, unlimitedMemMonitor,
	)
	limitedMemMonitor.Start(ctx, unlimitedMemMonitor, mon.BoundAccount{})
	defer limitedMemMonitor.Stop(ctx)
	limitedMemAcc := limitedMemMonitor.MakeBoundAccount()
	defer limitedMemAcc.Close(ctx)

	evalCtx := tree.MakeTestingEvalContext(st)
	testColumnFactory := coldataext.NewExtendedColumnFactory(&evalCtx)
	allocator := NewLimitedAllocator(ctx, &limitedMemAcc, &unlimitedMemAcc, testColumnFactory)

	// Check that no error occurs if the limit is not exceeded.
	require.NotPanics(t, func() { allocator.AdjustMemoryUsage(limit / 2) })
	require.Equal(t, limit/2, limitedMemAcc.Used())
	require.Zero(t, unlimitedMemAcc.Used())

	// Exceed the limit "before" making an allocation and ensure that the
	// unlimited account has not been grown.
	err := colexecerror.CatchVectorizedRuntimeError(func() { allocator.AdjustMemoryUsage(limit) })
	require.NotNil(t, err)
	require.True(t, strings.Contains(err.Error(), limitedMemMonitorName))
	require.Equal(t, limit/2, limitedMemAcc.Used())
	require.Zero(t, unlimitedMemAcc.Used())

	// Now exceed the limit "after" making an allocation and ensure that the
	// unlimited account has been grown.
	err = colexecerror.CatchVectorizedRuntimeError(func() { allocator.adjustMemoryUsageAfterAllocation(limit) })
	require.NotNil(t, err)
	require.True(t, strings.Contains(err.Error(), limitedMemMonitorName))
	require.Equal(t, limit/2, limitedMemAcc.Used())
	require.Equal(t, limit, unlimitedMemAcc.Used())

	// Ensure that the error from the unlimited memory account is returned when
	// it cannot be grown.
	err = colexecerror.CatchVectorizedRuntimeError(func() { allocator.adjustMemoryUsageAfterAllocation(math.MaxInt64) })
	require.NotNil(t, err)
	require.False(t, strings.Contains(err.Error(), limitedMemMonitorName))
	require.Equal(t, limit/2, limitedMemAcc.Used())
	require.Equal(t, limit, unlimitedMemAcc.Used())
}
