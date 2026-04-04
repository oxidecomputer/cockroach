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

	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats/ssmemstorage"
)

// Sink provides clients with interfaces to send statistics data into the sink.
type Sink interface {
	// AddAppStats ingests a single ssmemstorage.Container for a given appName.
	AddAppStats(ctx context.Context, appName string, other *ssmemstorage.Container) error
}

var _ Sink = &SQLStats{}

// AddAppStats implements the Sink interface.
func (s *SQLStats) AddAppStats(
	ctx context.Context, appName string, other *ssmemstorage.Container,
) error {
	stats := s.getStatsForApplication(appName)
	// Container.Add() manages locks for itself, so we don't need to guard it
	// with locks.
	return stats.Add(ctx, other)
}
