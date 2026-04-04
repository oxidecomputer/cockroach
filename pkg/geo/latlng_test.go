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

package geo

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeLongitudeDegrees(t *testing.T) {
	testCases := []struct {
		lng      float64
		expected float64
	}{
		{180, 180},
		{-180, -180},
		{181, -179},
		{360, 0},
		{-360, 0},
		{95, 95},
		{0, 0},
		{-10, -10},
		{10, 10},
		{555, -165},
		{-555, 165},
	}

	for _, tc := range testCases {
		t.Run(strconv.FormatFloat(tc.lng, 'f', -1, 64), func(t *testing.T) {
			require.Equal(t, tc.expected, NormalizeLongitudeDegrees(tc.lng))
		})
	}
}

func TestNormalizeLatitudeDegrees(t *testing.T) {
	testCases := []struct {
		lat      float64
		expected float64
	}{
		{0, 0},
		{10, 10},
		{-10, -10},
		{95, 85},
		{-95, -85},
		{90, 90},
		{-90, -90},
		{-180, 0},
		{180, 0},
		{270, -90},
		{-270, 90},
		{555, -15},
		{-555, 15},
	}

	for _, tc := range testCases {
		t.Run(strconv.FormatFloat(tc.lat, 'f', -1, 64), func(t *testing.T) {
			require.Equal(t, tc.expected, NormalizeLatitudeDegrees(tc.lat))
		})
	}
}
