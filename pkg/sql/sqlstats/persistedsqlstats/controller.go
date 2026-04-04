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

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats/sslocal"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlutil"
)

// Controller implements the SQL Stats subsystem control plane. This exposes
// administrative interfaces that can be consumed by other parts of the database
// (e.g. status server, builtins) to control the behavior of the SQL Stats
// subsystem.
type Controller struct {
	*sslocal.Controller
	db *kv.DB
	ie sqlutil.InternalExecutor
	st *cluster.Settings
}

// NewController returns a new instance of sqlstats.Controller.
func NewController(
	sqlStats *PersistedSQLStats,
	status serverpb.SQLStatusServer,
	db *kv.DB,
	ie sqlutil.InternalExecutor,
) *Controller {
	return &Controller{
		Controller: sslocal.NewController(sqlStats.SQLStats, status),
		db:         db,
		ie:         ie,
		st:         sqlStats.cfg.Settings,
	}
}

// CreateSQLStatsCompactionSchedule implements the tree.SQLStatsController
// interface.
func (s *Controller) CreateSQLStatsCompactionSchedule(ctx context.Context) error {
	return s.db.Txn(ctx, func(ctx context.Context, txn *kv.Txn) error {
		_, err := CreateSQLStatsCompactionScheduleIfNotYetExist(ctx, s.ie, txn, s.st)
		return err
	})
}

// ResetClusterSQLStats implements the tree.SQLStatsController interface. This
// method resets both the cluster-wide in-memory stats (via RPC fanout) and
// persisted stats (via TRUNCATE SQL statement)
func (s *Controller) ResetClusterSQLStats(ctx context.Context) error {
	if err := s.Controller.ResetClusterSQLStats(ctx); err != nil {
		return err
	}

	resetSysTableStats := func(tableName string) error {
		if _, err := s.ie.ExecEx(
			ctx,
			"reset-sql-stats",
			nil, /* txn */
			sessiondata.InternalExecutorOverride{
				User: security.NodeUserName(),
			},
			"TRUNCATE "+tableName); err != nil {
			return err
		}

		return nil
	}
	if err := resetSysTableStats("system.statement_statistics"); err != nil {
		return err
	}

	return resetSysTableStats("system.transaction_statistics")
}
