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

package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/stretchr/testify/require"
)

func TestCreateScheduledJob(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	h, cleanup := newTestHelper(t)
	defer cleanup()

	j := h.newScheduledJob(t, "test_job", "test sql")
	require.NoError(t, j.SetSchedule("@daily"))
	require.NoError(t, j.Create(context.Background(), h.cfg.InternalExecutor, nil))
	require.True(t, j.ScheduleID() > 0)
}

func TestCreatePausedScheduledJob(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	h, cleanup := newTestHelper(t)
	defer cleanup()

	j := h.newScheduledJob(t, "test_job", "test sql")
	require.NoError(t, j.SetSchedule("@daily"))
	j.Pause()
	require.NoError(t, j.Create(context.Background(), h.cfg.InternalExecutor, nil))
	require.True(t, j.ScheduleID() > 0)
	require.True(t, j.NextRun().Equal(time.Time{}))
}

func TestSetsSchedule(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	h, cleanup := newTestHelper(t)
	defer cleanup()

	j := h.newScheduledJob(t, "test_job", "test sql")

	// Set job schedule to run "@daily" -- i.e. at midnight.
	require.NoError(t, j.SetSchedule("@daily"))

	// The job is expected to run at midnight the next day.
	// We want to ensure nextRun correctly persisted in the cron table.
	expectedNextRun := h.env.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)

	require.NoError(t, j.Create(context.Background(), h.cfg.InternalExecutor, nil))

	loaded := h.loadSchedule(t, j.ScheduleID())
	require.Equal(t, j.ScheduleID(), loaded.ScheduleID())
	require.Equal(t, "@daily", loaded.rec.ScheduleExpr)
	require.True(t, loaded.NextRun().Equal(expectedNextRun))
}

func TestCreateOneOffJob(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	h, cleanup := newTestHelper(t)
	defer cleanup()

	j := h.newScheduledJob(t, "test_job", "test sql")
	j.SetNextRun(timeutil.Now())

	require.NoError(t, j.Create(context.Background(), h.cfg.InternalExecutor, nil))
	require.True(t, j.ScheduleID() > 0)

	loaded := h.loadSchedule(t, j.ScheduleID())
	require.Equal(t, j.NextRun().Round(time.Microsecond), loaded.NextRun())
	require.Equal(t, "", loaded.rec.ScheduleExpr)
}

func TestPauseUnpauseJob(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	h, cleanup := newTestHelper(t)
	defer cleanup()

	ctx := context.Background()
	j := h.newScheduledJob(t, "test_job", "test sql")
	require.NoError(t, j.SetSchedule("@daily"))
	require.NoError(t, j.Create(ctx, h.cfg.InternalExecutor, nil))

	// Pause and save.
	j.Pause()
	require.NoError(t, j.Update(ctx, h.cfg.InternalExecutor, nil))

	// Verify job is paused
	loaded := h.loadSchedule(t, j.ScheduleID())
	// Paused jobs have next run time set to NULL
	require.True(t, loaded.IsPaused())

	// Un-pausing the job resets next run time.
	require.NoError(t, j.ScheduleNextRun())
	require.NoError(t, j.Update(ctx, h.cfg.InternalExecutor, nil))

	// Verify job is no longer paused
	loaded = h.loadSchedule(t, j.ScheduleID())
	// Running schedules have nextRun set to non-null value
	require.False(t, loaded.IsPaused())
	require.False(t, loaded.NextRun().Equal(time.Time{}))
}
