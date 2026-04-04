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

package persistedsqlstats

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats/ssmemstorage"
	"github.com/cockroachdb/errors"
)

// ApplicationStats is a sqlstats.ApplicationStats that wraps an in-memory
// node-local ApplicationStats. ApplicationStats signals the subsystem when it
// encounters memory pressure which will triggers the flush operation.
type ApplicationStats struct {
	// local in-memory storage.
	sqlstats.ApplicationStats

	// Use to signal the stats writer is experiencing memory pressure.
	memoryPressureSignal chan struct{}
}

var _ sqlstats.ApplicationStats = &ApplicationStats{}

// RecordStatement implements sqlstats.ApplicationStats interface.
func (s *ApplicationStats) RecordStatement(
	ctx context.Context, key roachpb.StatementStatisticsKey, value sqlstats.RecordedStmtStats,
) (roachpb.StmtFingerprintID, error) {
	var fingerprintID roachpb.StmtFingerprintID
	err := s.recordStatsOrSendMemoryPressureSignal(func() (err error) {
		fingerprintID, err = s.ApplicationStats.RecordStatement(ctx, key, value)
		return err
	})
	return fingerprintID, err
}

// ShouldSaveLogicalPlanDesc implements sqlstats.ApplicationStats interface.
func (s *ApplicationStats) ShouldSaveLogicalPlanDesc(
	fingerprint string, implicitTxn bool, database string,
) bool {
	return s.ApplicationStats.ShouldSaveLogicalPlanDesc(fingerprint, implicitTxn, database)
}

// RecordTransaction implements sqlstats.ApplicationStats interface and saves
// per-transaction statistics.
func (s *ApplicationStats) RecordTransaction(
	ctx context.Context, key roachpb.TransactionFingerprintID, value sqlstats.RecordedTxnStats,
) error {
	return s.recordStatsOrSendMemoryPressureSignal(func() error {
		return s.ApplicationStats.RecordTransaction(ctx, key, value)
	})
}

func (s *ApplicationStats) recordStatsOrSendMemoryPressureSignal(fn func() error) error {
	err := fn()
	if errors.Is(err, ssmemstorage.ErrFingerprintLimitReached) || errors.Is(err, ssmemstorage.ErrMemoryPressure) {
		select {
		case s.memoryPressureSignal <- struct{}{}:
			// If we successfully signaled that we are experiencing memory pressure,
			// then our job is done. However, if we fail to send the signal, that
			// means we are already experiencing memory pressure and the
			// stats-flush-worker has already started to handle the flushing. We
			// don't need to do anything here at this point. The default case of the
			// select allows this operation to be non-blocking.
		default:
		}
		// We have already handled the memory pressure error. We don't have to
		// bubble up the error any further.
		return nil
	}
	return err
}
