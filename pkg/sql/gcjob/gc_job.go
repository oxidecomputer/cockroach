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

package gcjob

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
)

var (
	// MaxSQLGCInterval is the longest the polling interval between checking if
	// elements should be GC'd.
	MaxSQLGCInterval = 5 * time.Minute
)

// SetSmallMaxGCIntervalForTest sets the MaxSQLGCInterval and then returns a closure
// that resets it.
// This is to be used in tests like:
//    defer SetSmallMaxGCIntervalForTest()
func SetSmallMaxGCIntervalForTest() func() {
	oldInterval := MaxSQLGCInterval
	MaxSQLGCInterval = 500 * time.Millisecond
	return func() {
		MaxSQLGCInterval = oldInterval
	}
}

var idleWaitDuration = settings.RegisterDurationSetting(
	settings.TenantReadOnly,
	"sql.gc_job.idle_wait_duration",
	"after this duration of waiting for an update, the gc job will mark itself idle",
	time.Second,
)

type schemaChangeGCResumer struct {
	job *jobs.Job
}

// performGC GCs any schema elements that are in the DELETING state and returns
// a bool indicating if it GC'd any elements.
func performGC(
	ctx context.Context,
	execCfg *sql.ExecutorConfig,
	details *jobspb.SchemaChangeGCDetails,
	progress *jobspb.SchemaChangeGCProgress,
) error {
	if details.Tenant != nil {
		return errors.Wrapf(
			gcTenant(ctx, execCfg, details.Tenant.ID, progress),
			"attempting to GC tenant %+v", details.Tenant,
		)
	}
	if details.Indexes != nil {
		return errors.Wrap(gcIndexes(ctx, execCfg, details.ParentID, progress), "attempting to GC indexes")
	} else if details.Tables != nil {
		if err := gcTables(ctx, execCfg, progress); err != nil {
			return errors.Wrap(err, "attempting to GC tables")
		}

		// Drop database zone config when all the tables have been GCed.
		if details.ParentID != descpb.InvalidID && isDoneGC(progress) {
			if err := deleteDatabaseZoneConfig(
				ctx,
				execCfg.DB,
				execCfg.Codec,
				execCfg.Settings,
				details.ParentID,
			); err != nil {
				return errors.Wrap(err, "deleting database zone config")
			}
		}
	}
	return nil
}

func unsplitRangesForTables(
	ctx context.Context,
	execCfg *sql.ExecutorConfig,
	droppedTables []jobspb.SchemaChangeGCDetails_DroppedID,
) error {
	if !execCfg.Codec.ForSystemTenant() {
		return nil
	}

	for _, droppedTable := range droppedTables {
		startKey := execCfg.Codec.TablePrefix(uint32(droppedTable.ID))
		span := roachpb.Span{
			Key:    startKey,
			EndKey: startKey.PrefixEnd(),
		}
		if err := sql.UnsplitRangesInSpan(ctx, execCfg.DB, span); err != nil {
			return err
		}
	}

	return nil
}

// unsplitRangesForIndexes unsplits ranges with dropped index in key prefix
func unsplitRangesForIndexes(
	ctx context.Context,
	execCfg *sql.ExecutorConfig,
	indexes []jobspb.SchemaChangeGCDetails_DroppedIndex,
	parentTableID descpb.ID,
) error {
	if !execCfg.Codec.ForSystemTenant() {
		return nil
	}

	for _, idx := range indexes {
		startKey := execCfg.Codec.IndexPrefix(uint32(parentTableID), uint32(idx.IndexID))
		idxSpan := roachpb.Span{
			Key:    startKey,
			EndKey: startKey.PrefixEnd(),
		}

		if err := sql.UnsplitRangesInSpan(ctx, execCfg.DB, idxSpan); err != nil {
			return err
		}
	}

	return nil
}

func maybeUnsplitRanges(
	ctx context.Context,
	execCfg *sql.ExecutorConfig,
	jobID jobspb.JobID,
	details *jobspb.SchemaChangeGCDetails,
	progress *jobspb.SchemaChangeGCProgress,
) error {
	if progress.RangesUnsplitDone {
		return nil
	}

	if len(details.Indexes) > 0 {
		if err := unsplitRangesForIndexes(ctx, execCfg, details.Indexes, details.ParentID); err != nil {
			return err
		}
	}

	if len(details.Tables) > 0 {
		if err := unsplitRangesForTables(ctx, execCfg, details.Tables); err != nil {
			return err
		}
	}

	progress.RangesUnsplitDone = true
	persistProgress(ctx, execCfg, jobID, progress, runningStatusGC(progress))

	return nil
}

// Resume is part of the jobs.Resumer interface.
func (r schemaChangeGCResumer) Resume(ctx context.Context, execCtx interface{}) (err error) {
	defer func() {
		if err != nil && !r.isPermanentGCError(err) {
			err = jobs.MarkAsRetryJobError(err)
		}
	}()
	p := execCtx.(sql.JobExecContext)
	// TODO(pbardea): Wait for no versions.

	// Clone the ExecConfig so that fields can be overwritten for testing knobs.
	execCfg := *p.ExecCfg()
	if n := execCfg.GCJobTestingKnobs.Notifier; n != nil {
		execCfg.GCJobNotifier = n
	}
	// Use the same SystemConfigProvider as the notifier.
	execCfg.SystemConfig = execCfg.GCJobNotifier.SystemConfigProvider()

	if err := execCfg.JobRegistry.CheckPausepoint("gcjob.before_resume"); err != nil {
		return err
	}

	if fn := execCfg.GCJobTestingKnobs.RunBeforeResume; fn != nil {
		if err := fn(r.job.ID()); err != nil {
			return err
		}
	}
	details, progress, err := initDetailsAndProgress(ctx, &execCfg, r.job.ID())
	if err != nil {
		return err
	}

	if err := maybeUnsplitRanges(ctx, &execCfg, r.job.ID(), details, progress); err != nil {
		return err
	}

	tableDropTimes, indexDropTimes := getDropTimes(details)

	gossipUpdateC, cleanup := execCfg.GCJobNotifier.AddNotifyee(ctx)
	defer cleanup()
	var timerDuration time.Duration
	ts := timeutil.DefaultTimeSource{}

	for {
		idleWait := idleWaitDuration.Get(execCfg.SV())
		if err := waitForWork(
			ctx, r.job.MarkIdle, ts, timerDuration, idleWait, gossipUpdateC,
		); err != nil {
			return err
		}

		// Refresh the status of all elements in case any GC TTLs have changed.
		var expired bool
		var earliestDeadline time.Time
		if details.Tenant == nil {
			remainingTables := getAllTablesWaitingForGC(details, progress)
			expired, earliestDeadline = refreshTables(
				ctx, &execCfg, remainingTables, tableDropTimes, indexDropTimes, r.job.ID(), progress,
			)
		} else {
			expired, earliestDeadline, err = refreshTenant(ctx, &execCfg, details.Tenant.DropTime, details, progress)
			if err != nil {
				return err
			}
		}
		timerDuration = time.Until(earliestDeadline)

		if expired {
			// Some elements have been marked as DELETING so save the progress.
			persistProgress(ctx, &execCfg, r.job.ID(), progress, runningStatusGC(progress))
			if fn := execCfg.GCJobTestingKnobs.RunBeforePerformGC; fn != nil {
				if err := fn(r.job.ID()); err != nil {
					return err
				}
			}
			if err := performGC(ctx, &execCfg, details, progress); err != nil {
				return err
			}
			persistProgress(ctx, &execCfg, r.job.ID(), progress, sql.RunningStatusWaitingGC)

			// Trigger immediate re-run in case of more expired elements.
			timerDuration = 0
		}

		if isDoneGC(progress) {
			return nil
		}

		// Schedule the next check for GC.
		if timerDuration > MaxSQLGCInterval {
			timerDuration = MaxSQLGCInterval
		}
	}
}

// waitForWork waits until there is work to do given the gossipUpDateC, the
// timer, or the context. It calls markIdle with true after waiting
// idleWaitDuration. It calls markIdle with false before returning.
func waitForWork(
	ctx context.Context,
	markIdle func(isIdle bool),
	source timeutil.TimeSource,
	workTimerDuration, idleWaitDuration time.Duration,
	gossipUpdateC <-chan struct{},
) error {
	var markedIdle bool
	defer func() {
		if markedIdle {
			markIdle(false)
		}
	}()

	markIdleTimer := source.NewTimer()
	markIdleTimer.Reset(idleWaitDuration)
	defer markIdleTimer.Stop()

	workTimer := source.NewTimer()
	workTimer.Reset(workTimerDuration)
	defer workTimer.Stop()

	wait := func() (done bool) {
		select {
		case <-markIdleTimer.Ch():
			markIdleTimer.MarkRead()
			markIdle(true)
			markedIdle = true
			return false

		case <-gossipUpdateC:
			if log.V(2) {
				log.Info(ctx, "received a new system config")
			}

		case <-workTimer.Ch():
			workTimer.MarkRead()
			if log.V(2) {
				log.Info(ctx, "SchemaChangeGC workTimer triggered")
			}

		case <-ctx.Done():
		}
		return true
	}
	if done := wait(); !done {
		wait()
	}
	return ctx.Err()
}

// OnFailOrCancel is part of the jobs.Resumer interface.
func (r schemaChangeGCResumer) OnFailOrCancel(context.Context, interface{}, error) error {
	return nil
}

// isPermanentGCError returns true if the error is a permanent job failure,
// which indicates that the failed GC job cannot be retried.
func (r *schemaChangeGCResumer) isPermanentGCError(err error) bool {
	// Currently we classify errors based on Schema Change function to backport
	// to 20.2 and 21.1. This functionality should be changed once #44594 is
	// implemented.
	return sql.IsPermanentSchemaChangeError(err)
}

func init() {
	createResumerFn := func(job *jobs.Job, settings *cluster.Settings) jobs.Resumer {
		return &schemaChangeGCResumer{
			job: job,
		}
	}
	jobs.RegisterConstructor(jobspb.TypeSchemaChangeGC, createResumerFn, jobs.UsesTenantCostControl)
}
