// Copyright 2017 The Cockroach Authors.
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

package jobsprotectedts

import (
	"context"
	"strconv"

	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts/ptpb"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/protectedts/ptreconcile"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/scheduledjobs"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlutil"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
	"github.com/cockroachdb/errors"
)

// MetaType represents the types of meta values we support for records
// associated with jobs or schedules.
type MetaType int

const (
	// Jobs is the meta type for records associated with jobs.
	Jobs MetaType = iota
	// Schedules is the meta type for records associated with schedules.
	Schedules
)

// The value of metaTypes is used in the ptpb.Record.MetaType field for records
// associated with jobs/schedules.
//
// These values must not be changed as it is used durably in the database.
var metaTypes = map[MetaType]string{Jobs: "jobs", Schedules: "schedules"}

// GetMetaType return the value for the provided metaType that is used in the
// ptpb.Record.MetaType field for records associated with jobs/schedules.
func GetMetaType(metaType MetaType) string {
	return metaTypes[metaType]
}

// MakeStatusFunc returns a function which determines whether the job or
// schedule implied with this value of meta should be removed by the reconciler.
func MakeStatusFunc(
	jr *jobs.Registry, ie sqlutil.InternalExecutor, metaType MetaType,
) ptreconcile.StatusFunc {
	switch metaType {
	case Jobs:
		return func(ctx context.Context, txn *kv.Txn, meta []byte) (shouldRemove bool, _ error) {
			jobID, err := decodeID(meta)
			if err != nil {
				return false, err
			}
			j, err := jr.LoadJobWithTxn(ctx, jobspb.JobID(jobID), txn)
			if jobs.HasJobNotFoundError(err) {
				return true, nil
			}
			if err != nil {
				return false, err
			}
			isTerminal := j.CheckTerminalStatus(ctx, txn)
			return isTerminal, nil
		}
	case Schedules:
		return func(ctx context.Context, txn *kv.Txn, meta []byte) (shouldRemove bool, _ error) {
			scheduleID, err := decodeID(meta)
			if err != nil {
				return false, err
			}
			_, err = jobs.LoadScheduledJob(ctx, scheduledjobs.ProdJobSchedulerEnv, scheduleID, ie, txn)
			if jobs.HasScheduledJobNotFoundError(err) {
				return true, nil
			}
			return false, err
		}
	}
	return nil
}

// MakeRecord makes a protected timestamp record to protect a timestamp on
// behalf of this job.
//
// TODO(adityamaru): In 22.2 stop passing `deprecatedSpans` since PTS records
// will stop protecting key spans.
func MakeRecord(
	recordID uuid.UUID,
	metaID int64,
	tsToProtect hlc.Timestamp,
	deprecatedSpans []roachpb.Span,
	metaType MetaType,
	target *ptpb.Target,
) *ptpb.Record {
	return &ptpb.Record{
		ID:              recordID.GetBytesMut(),
		Timestamp:       tsToProtect,
		Mode:            ptpb.PROTECT_AFTER,
		MetaType:        metaTypes[metaType],
		Meta:            encodeID(metaID),
		DeprecatedSpans: deprecatedSpans,
		Target:          target,
	}
}

func encodeID(id int64) []byte {
	return []byte(strconv.FormatInt(id, 10))
}

func decodeID(meta []byte) (id int64, err error) {
	id, err = strconv.ParseInt(string(meta), 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to interpret meta %q as bytes", meta)
	}
	return id, err
}
