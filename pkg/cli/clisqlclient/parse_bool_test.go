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

package clisqlclient

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBool(t *testing.T) {
	defer leaktest.AfterTest(t)()
	testcases := []struct {
		input     string
		expect    bool
		expectErr bool
	}{
		{"true", true, false},
		{"on", true, false},
		{"yes", true, false},
		{"1", true, false},
		{" TrUe	", true, false},

		{"false", false, false},
		{"off", false, false},
		{"no", false, false},
		{"0", false, false},
		{"	FaLsE ", false, false},

		{"", false, true},
		{"foo", false, true},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			b, err := ParseBool(tc.input)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expect, b)
			}
		})
	}
}
