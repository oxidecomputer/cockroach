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
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats/sslocal"
)

// memStmtStatsIterator wraps a sslocal.StmtStatsIterator. Since in-memory
// statement statistics does not have aggregated_ts and aggregation_interval
// fields populated, memStmtStatsIterator overrides the
// sslocal.StmtStatsIterator's Cur() method to populate the aggregated_ts
// and aggregation_interval fields on the returning
// roachpb.CollectedStatementStatistics.
type memStmtStatsIterator struct {
	*sslocal.StmtStatsIterator
	aggregatedTs time.Time
	aggInterval  time.Duration
}

func newMemStmtStatsIterator(
	stats *sslocal.SQLStats,
	options *sqlstats.IteratorOptions,
	aggregatedTS time.Time,
	aggInterval time.Duration,
) *memStmtStatsIterator {
	return &memStmtStatsIterator{
		StmtStatsIterator: stats.StmtStatsIterator(options),
		aggregatedTs:      aggregatedTS,
		aggInterval:       aggInterval,
	}
}

// Cur calls the m.StmtStatsIterator.Cur() and populates the c.AggregatedTs
// field and c.AggregationInterval field.
func (m *memStmtStatsIterator) Cur() *roachpb.CollectedStatementStatistics {
	c := m.StmtStatsIterator.Cur()
	c.AggregatedTs = m.aggregatedTs
	c.AggregationInterval = m.aggInterval
	return c
}

// memTxnStatsIterator wraps a sslocal.TxnStatsIterator. Since in-memory
// transaction statistics does not have aggregated_ts and aggregation_interval
// fields populated, memTxnStatsIterator overrides the
// sslocal.TxnStatsIterator's Cur() method to populate the aggregated_ts and
// aggregatoin_interval fields fields on the returning
// roachpb.CollectedTransactionStatistics.
type memTxnStatsIterator struct {
	*sslocal.TxnStatsIterator
	aggregatedTs time.Time
	aggInterval  time.Duration
}

func newMemTxnStatsIterator(
	stats *sslocal.SQLStats,
	options *sqlstats.IteratorOptions,
	aggregatedTS time.Time,
	aggInterval time.Duration,
) *memTxnStatsIterator {
	return &memTxnStatsIterator{
		TxnStatsIterator: stats.TxnStatsIterator(options),
		aggregatedTs:     aggregatedTS,
		aggInterval:      aggInterval,
	}
}

// Cur calls the m.TxnStatsIterator.Cur() and populates the stats.AggregatedTs
// and stats.AggregationInterval fields.
func (m *memTxnStatsIterator) Cur() *roachpb.CollectedTransactionStatistics {
	stats := m.TxnStatsIterator.Cur()
	stats.AggregatedTs = m.aggregatedTs
	stats.AggregationInterval = m.aggInterval
	return stats
}
