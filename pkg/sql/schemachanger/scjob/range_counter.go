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

package scjob

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scdeps"
)

// rangeCounter implements scdeps.RangeCounter
type rangeCounter struct {
	db  *kv.DB
	dsp *sql.DistSQLPlanner
}

// NewRangeCounter constructs a new RangeCounter.
func NewRangeCounter(db *kv.DB, dsp *sql.DistSQLPlanner) scdeps.RangeCounter {
	return &rangeCounter{
		db:  db,
		dsp: dsp,
	}
}

var _ scdeps.RangeCounter = (*rangeCounter)(nil)

func (r rangeCounter) NumRangesInSpanContainedBy(
	ctx context.Context, span roachpb.Span, containedBy []roachpb.Span,
) (total, inContainedBy int, _ error) {
	return sql.NumRangesInSpanContainedBy(ctx, r.db, r.dsp, span, containedBy)
}
