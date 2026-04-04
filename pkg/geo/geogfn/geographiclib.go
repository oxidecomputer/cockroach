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
	"github.com/cockroachdb/cockroach/pkg/geo/geographiclib"
	"github.com/golang/geo/s2"
)

// spheroidDistance returns the s12 (meter) component of spheroid.Inverse from s2 Points.
func spheroidDistance(s *geographiclib.Spheroid, a s2.Point, b s2.Point) float64 {
	inv, _, _ := s.Inverse(s2.LatLngFromPoint(a), s2.LatLngFromPoint(b))
	return inv
}
