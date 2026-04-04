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

package colexecutils

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/colexec/colexectestutils"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestBoolVecToSelOp(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	tcs := []struct {
		boolCol  uint32
		tuples   colexectestutils.Tuples
		expected colexectestutils.Tuples
	}{
		{
			boolCol:  0,
			tuples:   colexectestutils.Tuples{{true}, {false}, {true}},
			expected: colexectestutils.Tuples{{true}, {true}},
		},
	}
	for _, tc := range tcs {
		colexectestutils.RunTests(t, testAllocator, []colexectestutils.Tuples{tc.tuples}, tc.expected, colexectestutils.OrderedVerifier, func(input []colexecop.Operator) (colexecop.Operator, error) {
			return NewBoolVecToSelOp(input[0], 0), nil
		})
	}
}
