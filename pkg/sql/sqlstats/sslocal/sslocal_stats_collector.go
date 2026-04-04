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

package sslocal

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/sessionphase"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// StatsCollector is used to collect statement and transaction statistics
// from connExecutor.
type StatsCollector struct {
	sqlstats.ApplicationStats

	// phaseTimes tracks session-level phase times.
	phaseTimes *sessionphase.Times

	// previousPhaseTimes tracks the session-level phase times for the previous
	// query. This enables the `SHOW LAST QUERY STATISTICS` observer statement.
	previousPhaseTimes *sessionphase.Times

	flushTarget sqlstats.ApplicationStats
	st          *cluster.Settings
	knobs       *sqlstats.TestingKnobs
}

var _ sqlstats.ApplicationStats = &StatsCollector{}

// NewStatsCollector returns an instance of sqlstats.StatsCollector.
func NewStatsCollector(
	st *cluster.Settings,
	appStats sqlstats.ApplicationStats,
	phaseTime *sessionphase.Times,
	knobs *sqlstats.TestingKnobs,
) *StatsCollector {
	return &StatsCollector{
		ApplicationStats: appStats,
		phaseTimes:       phaseTime.Clone(),
		st:               st,
		knobs:            knobs,
	}
}

// PhaseTimes implements sqlstats.StatsCollector interface.
func (s *StatsCollector) PhaseTimes() *sessionphase.Times {
	return s.phaseTimes
}

// PreviousPhaseTimes implements sqlstats.StatsCollector interface.
func (s *StatsCollector) PreviousPhaseTimes() *sessionphase.Times {
	return s.previousPhaseTimes
}

// Reset implements sqlstats.StatsCollector interface.
func (s *StatsCollector) Reset(appStats sqlstats.ApplicationStats, phaseTime *sessionphase.Times) {
	previousPhaseTime := s.phaseTimes
	s.flushTarget = appStats

	s.previousPhaseTimes = previousPhaseTime
	s.phaseTimes = phaseTime.Clone()
}

// StartTransaction implements sqlstats.StatsCollector interface.
// The current application stats are reset for the new transaction.
func (s *StatsCollector) StartTransaction() {
	s.flushTarget = s.ApplicationStats
	s.ApplicationStats = s.flushTarget.NewApplicationStatsWithInheritedOptions()
}

// EndTransaction implements sqlstats.StatsCollector interface.
func (s *StatsCollector) EndTransaction(
	ctx context.Context, transactionFingerprintID roachpb.TransactionFingerprintID,
) {
	// We possibly ignore the transactionFingerprintID, for situations where
	// grouping by it would otherwise result in collecting higher-cardinality
	// data in the system tables than the cleanup job is able to keep up with.
	// See #78338.
	if !AssociateStmtWithTxnFingerprint.Get(&s.st.SV) {
		transactionFingerprintID = roachpb.InvalidTransactionFingerprintID
	}

	var discardedStats uint64
	discardedStats += s.flushTarget.MergeApplicationStatementStats(
		ctx,
		s.ApplicationStats,
		func(statistics *roachpb.CollectedStatementStatistics) {
			statistics.Key.TransactionFingerprintID = transactionFingerprintID
		},
	)

	discardedStats += s.flushTarget.MergeApplicationTransactionStats(
		ctx,
		s.ApplicationStats,
	)

	if discardedStats > 0 {
		log.Warningf(ctx, "%d statement statistics discarded due to memory limit", discardedStats)
	}

	s.ApplicationStats.Free(ctx)
	s.ApplicationStats = s.flushTarget
	s.flushTarget = nil
}

// ShouldSaveLogicalPlanDesc implements sqlstats.StatsCollector interface.
func (s *StatsCollector) ShouldSaveLogicalPlanDesc(
	fingerprint string, implicitTxn bool, database string,
) bool {
	foundInFlushTarget := true

	if s.flushTarget != nil {
		foundInFlushTarget = s.flushTarget.ShouldSaveLogicalPlanDesc(fingerprint, implicitTxn, database)
	}

	return foundInFlushTarget &&
		s.ApplicationStats.ShouldSaveLogicalPlanDesc(fingerprint, implicitTxn, database)
}

// UpgradeImplicitTxn implements sqlstats.StatsCollector interface.
func (s *StatsCollector) UpgradeImplicitTxn(ctx context.Context) error {
	err := s.ApplicationStats.IterateStatementStats(ctx, &sqlstats.IteratorOptions{},
		func(_ context.Context, statistics *roachpb.CollectedStatementStatistics) error {
			statistics.Key.ImplicitTxn = false
			return nil
		})

	return err
}
