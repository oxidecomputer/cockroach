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

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCombinesUniqueInt64(t *testing.T) {
	for _, tc := range []struct{ inputA, inputB, expected []int64 }{
		{
			inputA:   []int64{1, 2, 4},
			inputB:   []int64{3, 5},
			expected: []int64{1, 2, 3, 4, 5},
		},
		{
			inputA:   []int64{1, 2, 4},
			inputB:   []int64{1, 3, 4},
			expected: []int64{1, 2, 3, 4},
		},
		{
			inputA:   []int64{1, 2, 3},
			inputB:   []int64{1, 2, 3},
			expected: []int64{1, 2, 3},
		},
		{
			inputA:   []int64{},
			inputB:   []int64{1, 3},
			expected: []int64{1, 3},
		},
	} {
		output := CombineUniqueInt64(tc.inputA, tc.inputB)
		require.Equal(t, tc.expected, output)
	}
}

func TestCombinesUniqueStrings(t *testing.T) {
	for _, tc := range []struct{ inputA, inputB, expected []string }{
		{
			inputA:   []string{"a", "b", "d"},
			inputB:   []string{"c", "e"},
			expected: []string{"a", "b", "c", "d", "e"},
		},
		{
			inputA:   []string{"a", "b", "d"},
			inputB:   []string{"a", "c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			inputA:   []string{"a", "b", "c"},
			inputB:   []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			inputA:   []string{},
			inputB:   []string{"a", "c"},
			expected: []string{"a", "c"},
		},
	} {
		output := CombineUniqueString(tc.inputA, tc.inputB)
		require.Equal(t, tc.expected, output)
	}
}
