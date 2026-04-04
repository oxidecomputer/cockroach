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

package rspb

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// FromTimestamp constructs a read summary from the provided timestamp, treating
// the argument as the low water mark of each segment in the summary.
func FromTimestamp(ts hlc.Timestamp) ReadSummary {
	seg := Segment{LowWater: ts}
	return ReadSummary{
		Local:  seg,
		Global: seg,
	}
}

// Clone performs a deep-copy of the receiver.
func (c ReadSummary) Clone() *ReadSummary {
	// NOTE: When ReadSummary is updated to include pointers to non-contiguous
	// memory, this will need to be updated.
	return &c
}

// Merge combines two read summaries, resulting in a single summary that
// reflects the combination of all reads in each original summary. The merge
// operation is commutative and idempotent.
func (c *ReadSummary) Merge(o ReadSummary) {
	c.Local.merge(o.Local)
	c.Global.merge(o.Global)
}

func (c *Segment) merge(o Segment) {
	c.LowWater.Forward(o.LowWater)
}

// AssertNoRegression asserts that all reads in the parameter's summary are
// reflected in the receiver's summary with at least as high of a timestamp.
func (c *ReadSummary) AssertNoRegression(ctx context.Context, o ReadSummary) {
	c.Local.assertNoRegression(ctx, o.Local, "local")
	c.Global.assertNoRegression(ctx, o.Global, "global")
}

func (c *Segment) assertNoRegression(ctx context.Context, o Segment, name string) {
	if c.LowWater.Less(o.LowWater) {
		log.Fatalf(ctx, "read summary regression in %s segment, was %s, now %s",
			name, o.LowWater, c.LowWater)
	}
}

// Ignore unused warning.
var _ = (*ReadSummary).AssertNoRegression
