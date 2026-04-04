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

package spanconfigtestutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpanRe(t *testing.T) {
	for _, tc := range []struct {
		input            string
		expMatch         bool
		expStart, expEnd string
	}{
		{"[a, b)", true, "a", "b"},
		{"[acd, bfg)", true, "acd", "bfg"}, // multi character keys allowed
		{"[a,b)", true, "a", "b"},          // separating space is optional
		{"[ a,b) ", false, "", ""},         // extraneous spaces disallowed
		{"[a,b ) ", false, "", ""},         // extraneous spaces disallowed
		{"[a,, b)", false, "", ""},         // only single comma allowed
		{" [a, b)", false, "", ""},         // need to start with '['
		{"[a,b)x", false, "", ""},          // need to end with ')'
	} {
		require.Equalf(t, tc.expMatch, spanRe.MatchString(tc.input), "input = %s", tc.input)
		if !tc.expMatch {
			continue
		}

		matches := spanRe.FindStringSubmatch(tc.input)
		require.Len(t, matches, 3)
		start, end := matches[1], matches[2]
		require.Equal(t, tc.expStart, start)
		require.Equal(t, tc.expEnd, end)
	}
}
