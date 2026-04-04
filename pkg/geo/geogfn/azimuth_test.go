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

package geogfn

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/geo"
	"github.com/stretchr/testify/require"
)

func TestAzimuth(t *testing.T) {
	testCases := []struct {
		desc     string
		a        string
		b        string
		expected float64
	}{
		{
			"north east",
			"POINT(0 0)",
			"POINT(1 1)",
			0.7886800845259658,
		},
		{
			"east",
			"POINT(0 0)",
			"POINT(1 0)",
			1.5707963267948966,
		},
		{
			"south east",
			"POINT(0 0)",
			"POINT(1 -1)",
			2.352912569063827,
		},
		{
			"srid set; south east",
			"SRID=4004;POINT(0 0)",
			"SRID=4004;POINT(1 -1)",
			2.3529226390237774,
		},
		{
			"due north",
			"POINT(0 0)",
			"POINT(0 1)",
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			a, err := geo.ParseGeography(tc.a)
			require.NoError(t, err)
			b, err := geo.ParseGeography(tc.b)
			require.NoError(t, err)

			r, err := Azimuth(a, b)
			require.NoError(t, err)
			require.NotNil(t, r)
			require.Equal(t, tc.expected, *r)
		})
	}

	t.Run("same point", func(t *testing.T) {
		a, err := geo.ParseGeography("POINT(1.0 1.0)")
		require.NoError(t, err)
		ret, err := Azimuth(a, a)
		require.NoError(t, err)
		require.Nil(t, ret)
	})

	errorTestCases := []struct {
		a          string
		b          string
		errorMatch string
	}{
		{
			"LINESTRING(0 0, 1 0)",
			"POINT(0 0)",
			"arguments must be POINT geometries",
		},
	}
	for _, tc := range errorTestCases {
		t.Run(fmt.Sprintf("%s <=> %s", tc.a, tc.b), func(t *testing.T) {
			a, err := geo.ParseGeography(tc.a)
			require.NoError(t, err)
			b, err := geo.ParseGeography(tc.b)
			require.NoError(t, err)

			_, err = Azimuth(a, b)
			require.Error(t, err)
			require.EqualError(t, err, tc.errorMatch)
		})
	}

	t.Run("errors if SRIDs mismatch", func(t *testing.T) {
		_, err := Azimuth(mismatchingSRIDGeographyA, mismatchingSRIDGeographyB)
		requireMismatchingSRIDError(t, err)
	})

	t.Run("errors on POINT EMPTY", func(t *testing.T) {
		_, err := Azimuth(
			geo.MustParseGeography("POINT EMPTY"),
			geo.MustParseGeography("POINT(1 0)"),
		)
		require.EqualError(t, err, "cannot call ST_Azimuth with POINT EMPTY")

		_, err = Azimuth(
			geo.MustParseGeography("POINT(1 0)"),
			geo.MustParseGeography("POINT EMPTY"),
		)
		require.EqualError(t, err, "cannot call ST_Azimuth with POINT EMPTY")
	})
}
