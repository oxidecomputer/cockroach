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

package props

import (
	"sort"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt"
)

func TestColumnStatisticsSort(t *testing.T) {
	type testCase struct {
		input    ColumnStatistics
		expected string
	}
	testCases := []testCase{
		{
			input: ColumnStatistics{
				{Cols: opt.MakeColSet(3)},
				{Cols: opt.MakeColSet(1)},
				{Cols: opt.MakeColSet(5)},
			},
			expected: "(1) (3) (5)",
		},
		{
			input: ColumnStatistics{
				{Cols: opt.MakeColSet(1, 3)},
				{Cols: opt.MakeColSet(1, 5)},
			},
			expected: "(1,3) (1,5)",
		},
		{
			input: ColumnStatistics{
				{Cols: opt.MakeColSet(3)},
				{Cols: opt.MakeColSet(1, 7)},
				{Cols: opt.MakeColSet(5)},
				{Cols: opt.MakeColSet(1, 3)},
				{Cols: opt.MakeColSet(1, 4, 6)},
				{Cols: opt.MakeColSet(1, 4, 7)},
			},
			expected: "(3) (5) (1,3) (1,7) (1,4,6) (1,4,7)",
		},
	}

	for _, tc := range testCases {
		sort.Sort(tc.input)
		var cols []string
		for i := 0; i < len(tc.input); i++ {
			cols = append(cols, tc.input[i].Cols.String())
		}
		result := strings.Join(cols, " ")
		if result != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, result)
		}
	}
}
