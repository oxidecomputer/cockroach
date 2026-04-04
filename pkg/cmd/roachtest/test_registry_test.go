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

package main

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/spec"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestMakeTestRegistry(t *testing.T) {
	testutils.RunTrueAndFalse(t, "preferSSD", func(t *testing.T, preferSSD bool) {
		r, err := makeTestRegistry(spec.AWS, "foo", "zone123", preferSSD)
		require.NoError(t, err)
		require.Equal(t, preferSSD, r.preferSSD)
		require.Equal(t, "zone123", r.zones)
		require.Equal(t, "foo", r.instanceType)
		require.Equal(t, spec.AWS, r.cloud)

		s := r.MakeClusterSpec(100, spec.Geo(), spec.Zones("zone99"), spec.CPU(12), spec.PreferSSD())
		require.EqualValues(t, 100, s.NodeCount)
		require.Equal(t, "foo", s.InstanceType)
		require.True(t, s.Geo)
		require.Equal(t, "zone99", s.Zones)
		require.EqualValues(t, 12, s.CPUs)
		require.True(t, s.PreferLocalSSD)
	})

}
