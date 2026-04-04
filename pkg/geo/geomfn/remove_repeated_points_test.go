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
	"github.com/cockroachdb/cockroach/pkg/geo/geopb"
	"github.com/stretchr/testify/require"
)

func TestRemoveRepeatedPoints(t *testing.T) {
	testCases := []struct {
		wkt      string
		expected string
	}{
		{"POINT EMPTY", "POINT EMPTY"},
		{"POINT (1 2)", "POINT (1 2)"},
		{"MULTIPOINT EMPTY", "MULTIPOINT EMPTY"},
		{"MULTIPOINT (1 2)", "MULTIPOINT (1 2)"},
		{"MULTIPOINT (1 2, 3 4)", "MULTIPOINT (1 2, 3 4)"},
		{"MULTIPOINT (1 2, 1 2, 3 4)", "MULTIPOINT (1 2, 3 4)"},
		{"MULTIPOINT (1 2, 3 4, 1 2)", "MULTIPOINT (1 2, 3 4, 1 2)"},
		{"MULTIPOINT (1 2, 1 2, 1 2, 1 2)", "MULTIPOINT (1 2)"},
		{"MULTIPOINT (1 2, EMPTY, 3 4)", "MULTIPOINT (1 2, EMPTY, 3 4)"},
		{"LINESTRING EMPTY", "LINESTRING EMPTY"},
		{"LINESTRING (1 2, 3 4)", "LINESTRING (1 2, 3 4)"},
		{"LINESTRING (1 2, 1 2, 3 4)", "LINESTRING (1 2, 3 4)"},
		{"LINESTRING (1 2, 3 4, 1 2)", "LINESTRING (1 2, 3 4, 1 2)"},
		{"LINESTRING (1 2, 1 2, 1 2, 1 2)", "LINESTRING (1 2, 1 2)"},
		{"MULTILINESTRING EMPTY", "MULTILINESTRING EMPTY"},
		{
			"MULTILINESTRING ((1 2, 3 4, 1 2), EMPTY, (1 2, 3 4))",
			"MULTILINESTRING ((1 2, 3 4, 1 2), EMPTY, (1 2, 3 4))",
		},
		{
			"MULTILINESTRING ((1 2, 1 2, 3 4), (1 2, 3 4, 3 4), (1 2, 1 2, 1 2, 1 2))",
			"MULTILINESTRING ((1 2, 3 4), (1 2, 3 4), (1 2, 1 2))",
		},
		{"POLYGON EMPTY", "POLYGON EMPTY"},
		{"POLYGON ((1 2, 3 4, 5 6, 1 2))", "POLYGON ((1 2, 3 4, 5 6, 1 2))"},
		{"POLYGON ((1 2, 3 4, 3 4, 5 6, 1 2))", "POLYGON ((1 2, 3 4, 5 6, 1 2))"},
		{
			"POLYGON ((1 2, 3 4, 5 6, 1 2), (1 2, 1 2, 3 4, 5 6, 1 2))",
			"POLYGON ((1 2, 3 4, 5 6, 1 2), (1 2, 3 4, 5 6, 1 2))",
		},
		{"POLYGON ((1 2, 1 2, 1 2, 1 2, 1 2, 1 2))", "POLYGON ((1 2, 1 2, 1 2, 1 2))"},
		{
			"MULTIPOLYGON (((1 2, 3 4, 5 6, 1 2)), EMPTY, ((1 2, 1 2, 3 4, 5 6, 1 2)))",
			"MULTIPOLYGON (((1 2, 3 4, 5 6, 1 2)), EMPTY, ((1 2, 3 4, 5 6, 1 2)))",
		},
		{"MULTIPOLYGON EMPTY", "MULTIPOLYGON EMPTY"},
		{
			"MULTIPOLYGON (((1 2, 3 4, 5 6, 1 2), (1 2, 1 2, 3 4, 5 6, 1 2)))",
			"MULTIPOLYGON (((1 2, 3 4, 5 6, 1 2), (1 2, 3 4, 5 6, 1 2)))",
		},
		{
			"MULTIPOLYGON (((1 2, 1 2, 1 2, 1 2, 1 2)), ((1 2, 1 2, 3 4, 5 6, 1 2)))",
			"MULTIPOLYGON (((1 2, 1 2, 1 2, 1 2)), ((1 2, 3 4, 5 6, 1 2)))",
		},
		{"GEOMETRYCOLLECTION EMPTY", "GEOMETRYCOLLECTION EMPTY"},
		{
			"GEOMETRYCOLLECTION (POINT (1 1), MULTIPOINT (1 2, 1 2, 3 4), LINESTRING (1 2, 3 4, 3 4))",
			"GEOMETRYCOLLECTION (POINT (1 1), MULTIPOINT (1 2, 3 4), LINESTRING (1 2, 3 4))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.wkt, func(t *testing.T) {
			srid := geopb.SRID(4000)
			g, err := geo.ParseGeometryFromEWKT(geopb.EWKT(tc.wkt), srid, true)
			require.NoError(t, err)

			res, err := RemoveRepeatedPoints(g, 0)
			require.NoError(t, err)
			wkt, err := geo.SpatialObjectToWKT(res.SpatialObject(), 0)
			require.NoError(t, err)
			require.EqualValues(t, tc.expected, wkt)
			require.EqualValues(t, srid, res.SRID())
		})
	}

	testCasesTolerance := []struct {
		wkt       string
		tolerance float64
		expected  string
	}{
		{"MULTIPOINT (1 1, 2 2, 3 3, 4 4, 5 5)", 1.4, "MULTIPOINT (1 1, 2 2, 3 3, 4 4, 5 5)"},
		{"MULTIPOINT (1 1, 2 2, 3 3, 4 4, 5 5)", 1.5, "MULTIPOINT (1 1, 3 3, 5 5)"},
		{"MULTIPOINT (1 1, 2 2, 3 3, 4 4, 5 5)", 3, "MULTIPOINT (1 1, 4 4)"},
		{"MULTIPOINT (1 1, 2 2, 3 3, 4 4, 5 5)", 6, "MULTIPOINT (1 1)"},
	}

	for _, tc := range testCasesTolerance {
		t.Run(fmt.Sprintf("%v tolerance %v", tc.wkt, tc.tolerance), func(t *testing.T) {
			srid := geopb.SRID(4000)
			g, err := geo.ParseGeometryFromEWKT(geopb.EWKT(tc.wkt), srid, true)
			require.NoError(t, err)

			res, err := RemoveRepeatedPoints(g, tc.tolerance)
			require.NoError(t, err)
			wkt, err := geo.SpatialObjectToWKT(res.SpatialObject(), 0)
			require.NoError(t, err)
			require.EqualValues(t, tc.expected, wkt)
			require.EqualValues(t, srid, res.SRID())
		})
	}
}
