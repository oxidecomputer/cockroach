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

package idxusage

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
)

// Controller implements the index usage stats subsystem control plane. This exposes
// administrative interfaces that can be consumed by other parts of the database
// (e.g. status server, builtins) to control the behavior of index usage stas
// subsystem.
type Controller struct {
	statusServer serverpb.SQLStatusServer
}

// NewController returns a new instance of idxusage.Controller.
func NewController(status serverpb.SQLStatusServer) *Controller {
	return &Controller{
		statusServer: status,
	}
}

// ResetIndexUsageStats implements the tree.IndexUsageStatsController interface.
func (s *Controller) ResetIndexUsageStats(ctx context.Context) error {
	req := &serverpb.ResetIndexUsageStatsRequest{}
	_, err := s.statusServer.ResetIndexUsageStats(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
