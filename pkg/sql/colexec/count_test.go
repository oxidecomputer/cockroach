// Copyright 2018 The Cockroach Authors.
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

package colexec

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexectestutils"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestCount(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	tcs := []struct {
		tuples   colexectestutils.Tuples
		expected colexectestutils.Tuples
	}{
		{
			tuples:   colexectestutils.Tuples{{1}, {1}},
			expected: colexectestutils.Tuples{{2}},
		},
		{
			tuples:   colexectestutils.Tuples{{1}, {1}, {1}, {1}, {1}, {1}, {1}, {1}, {1}, {1}},
			expected: colexectestutils.Tuples{{10}},
		},
	}
	for _, tc := range tcs {
		// The tuples consisting of all nulls still count as separate rows, so if
		// we replace all values with nulls, we should get the same output.
		colexectestutils.RunTestsWithoutAllNullsInjection(t, testAllocator, []colexectestutils.Tuples{tc.tuples}, nil, tc.expected, colexectestutils.OrderedVerifier, func(input []colexecop.Operator) (colexecop.Operator, error) {
			return NewCountOp(testAllocator, input[0]), nil
		})
	}
}
