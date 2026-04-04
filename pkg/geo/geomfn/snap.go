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

package geomfn

import (
	"github.com/cockroachdb/cockroach/pkg/geo"
	"github.com/cockroachdb/cockroach/pkg/geo/geos"
)

// Snap returns the input geometry with the vertices snapped to the target
// geometry. Tolerance is used to control where snapping is performed.
// If no snapping occurs then the input geometry is returned unchanged.
func Snap(input, target geo.Geometry, tolerance float64) (geo.Geometry, error) {
	snappedEWKB, err := geos.Snap(input.EWKB(), target.EWKB(), tolerance)
	if err != nil {
		return geo.Geometry{}, err
	}
	return geo.ParseGeometryFromEWKB(snappedEWKB)
}
