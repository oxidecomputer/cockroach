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

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestInlineExecutorFailedJobsHandling(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	argsFn := func(args *base.TestServerArgs) {
		args.Knobs.JobsTestingKnobs = NewTestingKnobsWithShortIntervals()
	}

	h, cleanup := newTestHelperWithServerArgs(t, argsFn)
	defer cleanup()

	var tests = []struct {
		onError         jobspb.ScheduleDetails_ErrorHandlingBehavior
		expectedNextRun time.Time
	}{
		{
			onError:         jobspb.ScheduleDetails_RETRY_SCHED,
			expectedNextRun: cronMustParse(t, "@daily").Next(h.env.Now()).Round(time.Microsecond),
		},
		{
			onError:         jobspb.ScheduleDetails_RETRY_SOON,
			expectedNextRun: h.env.Now().Add(retryFailedJobAfter).Round(time.Microsecond),
		},
		{
			onError:         jobspb.ScheduleDetails_PAUSE_SCHED,
			expectedNextRun: time.Time{}.UTC(),
		},
	}

	for _, test := range tests {
		t.Run(test.onError.String(), func(t *testing.T) {
			j := h.newScheduledJob(t, "test_job", "test sql")
			j.rec.ExecutorType = InlineExecutorName

			require.NoError(t, j.SetSchedule("@daily"))
			j.SetScheduleDetails(jobspb.ScheduleDetails{OnError: test.onError})

			ctx := context.Background()
			require.NoError(t, j.Create(ctx, h.cfg.InternalExecutor, nil))

			// Pretend we failed running; we expect job to be rescheduled.
			require.NoError(t, NotifyJobTermination(
				ctx, h.env, 123, StatusFailed, nil, j.ScheduleID(), h.cfg.InternalExecutor, nil))

			// Verify nextRun updated
			loaded := h.loadSchedule(t, j.ScheduleID())
			require.Equal(t, test.expectedNextRun, loaded.NextRun())
		})
	}
}
