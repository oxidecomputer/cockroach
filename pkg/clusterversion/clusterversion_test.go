// Copyright 2017 The Cockroach Authors.
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

package clusterversion

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/stretchr/testify/require"
)

// Test that the OnChange callback is called for the cluster version with the
// right argument.
func TestClusterVersionOnChange(t *testing.T) {
	ctx := context.Background()
	var sv settings.Values

	cvs := &clusterVersionSetting{}
	cvs.VersionSetting = settings.MakeVersionSetting(cvs)
	settings.RegisterVersionSetting(
		settings.TenantWritable,
		"dummy version key",
		"test description",
		&cvs.VersionSetting)

	handle := newHandleImpl(cvs, &sv, binaryVersion, binaryMinSupportedVersion)
	newCV := ClusterVersion{
		Version: roachpb.Version{
			Major:    1,
			Minor:    2,
			Patch:    3,
			Internal: 4,
		},
	}
	encoded, err := protoutil.Marshal(&newCV)
	require.NoError(t, err)

	var capturedV ClusterVersion
	handle.SetOnChange(func(ctx context.Context, newVersion ClusterVersion) {
		capturedV = newVersion
	})
	cvs.SetInternal(ctx, &sv, encoded)
	require.Equal(t, newCV, capturedV)
}
