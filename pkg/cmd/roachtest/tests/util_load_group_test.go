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

package tests

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/stretchr/testify/require"
)

func TestLoadGroups(t *testing.T) {
	for _, tc := range []struct {
		numZones, numRoachNodes, numLoadNodes int
		loadGroups                            loadGroupList
	}{
		{
			3, 9, 3,
			loadGroupList{
				{
					option.NodeListOption{1, 2, 3},
					option.NodeListOption{4},
				},
				{
					option.NodeListOption{5, 6, 7},
					option.NodeListOption{8},
				},
				{
					option.NodeListOption{9, 10, 11},
					option.NodeListOption{12},
				},
			},
		},
		{
			3, 9, 1,
			loadGroupList{
				{
					option.NodeListOption{1, 2, 3, 4, 5, 6, 7, 8, 9},
					option.NodeListOption{10},
				},
			},
		},
		{
			4, 8, 2,
			loadGroupList{
				{
					option.NodeListOption{1, 2, 3, 4},
					option.NodeListOption{9},
				},
				{
					option.NodeListOption{5, 6, 7, 8},
					option.NodeListOption{10},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("%d/%d/%d", tc.numZones, tc.numRoachNodes, tc.numLoadNodes),
			func(t *testing.T) {
				l := option.NodeLister{NodeCount: tc.numRoachNodes + tc.numLoadNodes, Fatalf: t.Fatalf}
				lg := makeLoadGroups(l, tc.numZones, tc.numRoachNodes, tc.numLoadNodes)
				require.EqualValues(t, lg, tc.loadGroups)
			})
	}
	t.Run("panics with too many load nodes", func(t *testing.T) {
		require.Panics(t, func() {

			numZones, numRoachNodes, numLoadNodes := 2, 4, 3
			makeLoadGroups(nil, numZones, numRoachNodes, numLoadNodes)
		}, "Failed to panic when number of load nodes exceeded number of zones")
	})
	t.Run("panics with unequal zones per load node", func(t *testing.T) {
		require.Panics(t, func() {
			numZones, numRoachNodes, numLoadNodes := 4, 4, 3
			makeLoadGroups(nil, numZones, numRoachNodes, numLoadNodes)
		}, "Failed to panic when number of zones is not divisible by number of load nodes")
	})
}
