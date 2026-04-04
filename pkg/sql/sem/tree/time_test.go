// Copyright 2019 The Cockroach Authors.
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

package tree

import (
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/assert"
)

func TestTimeFamilyPrecisionToRoundDuration(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	testCases := []struct {
		precision     int32
		expected      time.Duration
		expectedPanic bool
	}{
		{precision: 0, expected: time.Duration(1000000000)},
		{precision: 1, expected: time.Duration(100000000)},
		{precision: 2, expected: time.Duration(10000000)},
		{precision: 3, expected: time.Duration(1000000)},
		{precision: 4, expected: time.Duration(100000)},
		{precision: 5, expected: time.Duration(10000)},
		{precision: 6, expected: time.Duration(1000)},

		{precision: -2, expectedPanic: true},
		{precision: 7, expectedPanic: true},
	}

	for _, tc := range testCases {
		if tc.expectedPanic {
			assert.Panics(t, func() { TimeFamilyPrecisionToRoundDuration(tc.precision) })
		} else {
			assert.Equal(t, tc.expected, TimeFamilyPrecisionToRoundDuration(tc.precision))
		}
	}
}
