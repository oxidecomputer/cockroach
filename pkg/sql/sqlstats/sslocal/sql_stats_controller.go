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

	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// Controller implements the SQL Stats subsystem control plane. This exposes
// administrative interfaces that can be consumed by other parts of the database
// (e.g. status server, builtins) to control the behavior of the SQL Stats
// subsystem.
type Controller struct {
	sqlStats     *SQLStats
	statusServer serverpb.SQLStatusServer
}

// NewController returns a new instance of sqlstats.Controller.
func NewController(sqlStats *SQLStats, status serverpb.SQLStatusServer) *Controller {
	return &Controller{
		sqlStats:     sqlStats,
		statusServer: status,
	}
}

// ResetClusterSQLStats implements the tree.SQLStatsController interface.
func (s *Controller) ResetClusterSQLStats(ctx context.Context) error {
	req := &serverpb.ResetSQLStatsRequest{
		// ResetPersistedStats field is explicitly set to false here since
		// sslocal.Controller should only be resetting the in-memory stats.
		// persistedsqlstats.Controller is responsible for resetting persisted
		// stats.
		ResetPersistedStats: false,
	}
	_, err := s.statusServer.ResetSQLStats(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// ResetLocalSQLStats resets the node-local sql stats.
func (s *Controller) ResetLocalSQLStats(ctx context.Context) {
	err := s.sqlStats.Reset(ctx)
	if err != nil {
		if log.V(1) {
			log.Warningf(ctx, "reported SQL stats memory limit has been exceeded, some fingerprints stats are discarded: %s", err)
		}
	}
}
