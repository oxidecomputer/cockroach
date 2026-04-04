// Copyright 2020 The Cockroach Authors.
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

package geomfn

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/geo"
	"github.com/stretchr/testify/require"
)

func TestRelate(t *testing.T) {
	testCases := []struct {
		a        geo.Geometry
		b        geo.Geometry
		expected string
	}{
		{leftRect, rightRect, "FF2F11212"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc:%d", i), func(t *testing.T) {
			ret, err := Relate(tc.a, tc.b)
			require.NoError(t, err)
			require.Equal(t, tc.expected, ret)
		})
	}
}

func TestRelateBoundaryNodeRule(t *testing.T) {
	testCases := []struct {
		a        geo.Geometry
		b        geo.Geometry
		bnr      int
		expected string
	}{
		{leftRect, rightRect, 1, "FF2F11212"},
		{leftRect, rightRect, 2, "FF2F11212"},
		{leftRect, rightRect, 3, "1F2F002F2"},
		{leftRect, rightRect, 4, "FF2F11212"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc:%d", i), func(t *testing.T) {
			ret, err := RelateBoundaryNodeRule(tc.a, tc.b, tc.bnr)
			require.NoError(t, err)
			require.Equal(t, tc.expected, ret)
		})
	}

	t.Run("errors on invalid BNR", func(t *testing.T) {
		_, err := RelateBoundaryNodeRule(leftRect, rightRect, 0)
		require.Error(t, err)
		_, err = RelateBoundaryNodeRule(leftRect, rightRect, 5)
		require.Error(t, err)
	})

	t.Run("errors if SRIDs mismatch", func(t *testing.T) {
		_, err := RelateBoundaryNodeRule(mismatchingSRIDGeometryA, mismatchingSRIDGeometryB, 1)
		requireMismatchingSRIDError(t, err)
	})
}

func TestMatchesDE9IM(t *testing.T) {
	testCases := []struct {
		str           string
		pattern       string
		expected      bool
		expectedError string
	}{
		{"", "T**FF*FF*", false, `relation "" should be of length 9`},
		{"TTTTTTTTT", "T**FF*FF*T", false, `pattern "T**FF*FF*T" should be of length 9`},
		{"TTTTTTTTT", "T**FF*FF*T", false, `pattern "T**FF*FF*T" should be of length 9`},
		{"000FFF000", "cTTFfFTTT", false, `unrecognized pattern character: c`},
		{"120FFF021", "TTTFfFTTT", true, ""},
		{"02FFFF000", "T**FfFTTT", true, ""},
		{"02FFFF000", "T2*FfFTTT", true, ""},
		{"01FFFF000", "01*FfFTTT", true, ""},
		{"020F1F010", "TTTFFFTtT", false, ""},
		{"020FFF0f0", "TTTFFFTtT", false, ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s has pattern %s", tc.str, tc.pattern), func(t *testing.T) {
			ret, err := MatchesDE9IM(tc.str, tc.pattern)
			if tc.expectedError == "" {
				require.NoError(t, err)
				require.Equal(t, tc.expected, ret)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}

	t.Run("errors if SRIDs mismatch", func(t *testing.T) {
		_, err := Relate(mismatchingSRIDGeometryA, mismatchingSRIDGeometryB)
		requireMismatchingSRIDError(t, err)
	})
}
