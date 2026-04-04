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

package sql

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/backfill"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/stretchr/testify/require"
)

// TestingDistIndexBackfill exposes the index backfill functionality for
// testing.
func (sc *SchemaChanger) TestingDistIndexBackfill(
	ctx context.Context,
	version descpb.DescriptorVersion,
	targetSpans []roachpb.Span,
	addedIndexes []descpb.IndexID,
	filter backfill.MutationFilter,
) error {
	s := &multiStageFractionScaler{initial: 0.0, stages: backfillStageFractions}
	err := sc.distIndexBackfill(ctx, version, targetSpans, addedIndexes, true, filter, s)
	return err
}

// SetJob sets the job.
func (sc *SchemaChanger) SetJob(job *jobs.Job) {
	sc.job = job
}

func TestCalculateSplitAtShards(t *testing.T) {
	defer leaktest.AfterTest(t)()

	testCases := []struct {
		testName    string
		maxSplit    int64
		bucketCount int32
		expected    []int64
	}{
		{
			testName:    "buckets_less_than_max_split",
			maxSplit:    8,
			bucketCount: 0,
			expected:    []int64{},
		},
		{
			testName:    "buckets_less_than_max_split",
			maxSplit:    8,
			bucketCount: 5,
			expected:    []int64{0, 1, 2, 3, 4},
		},
		{
			testName:    "buckets_equal_to_max_split",
			maxSplit:    8,
			bucketCount: 8,
			expected:    []int64{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			testName:    "buckets_greater_than_max_split_1",
			maxSplit:    8,
			bucketCount: 30,
			expected:    []int64{0, 3, 7, 11, 15, 18, 22, 26},
		},
		{
			testName:    "buckets_greater_than_max_split_2",
			maxSplit:    8,
			bucketCount: 1000,
			expected:    []int64{0, 125, 250, 375, 500, 625, 750, 875},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			shards := calculateSplitAtShards(tc.maxSplit, tc.bucketCount)
			require.Equal(t, tc.expected, shards)
		})
	}
}
