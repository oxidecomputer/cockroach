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

package reports

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// computeConstraintConformanceReport iterates through all the ranges and
// generates the constraint conformance report.
func computeConstraintConformanceReport(
	ctx context.Context,
	rangeStore RangeIterator,
	cfg *config.SystemConfig,
	storeResolver StoreResolver,
) (ConstraintReport, error) {
	v := makeConstraintConformanceVisitor(ctx, cfg, storeResolver)
	err := visitRanges(ctx, rangeStore, cfg, &v)
	return v.Report(), err
}

// computeReplicationStatsReport iterates through all the ranges and generates
// the replication stats report.
func computeReplicationStatsReport(
	ctx context.Context, rangeStore RangeIterator, checker nodeChecker, cfg *config.SystemConfig,
) (RangeReport, error) {
	v := makeReplicationStatsVisitor(ctx, cfg, checker)
	err := visitRanges(ctx, rangeStore, cfg, &v)
	return v.Report(), err
}

// computeCriticalLocalitiesReport iterates through all the ranges and generates
// the critical localities report.
func computeCriticalLocalitiesReport(
	ctx context.Context,
	nodeLocalities map[roachpb.NodeID]roachpb.Locality,
	rangeStore RangeIterator,
	checker nodeChecker,
	cfg *config.SystemConfig,
	storeResolver StoreResolver,
) (LocalityReport, error) {
	v := makeCriticalLocalitiesVisitor(ctx, nodeLocalities, cfg, storeResolver, checker)
	err := visitRanges(ctx, rangeStore, cfg, &v)
	return v.Report(), err
}
