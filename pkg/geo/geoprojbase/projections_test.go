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

package geoprojbase

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjections(t *testing.T) {
	for srid, proj := range getProjections() {
		t.Run(strconv.Itoa(int(srid)), func(t *testing.T) {
			require.NotEqual(t, Bounds{}, proj.Bounds)
			require.GreaterOrEqual(t, proj.Bounds.MaxX, proj.Bounds.MinX)
			require.GreaterOrEqual(t, proj.Bounds.MaxY, proj.Bounds.MinY)
		})
	}
}
