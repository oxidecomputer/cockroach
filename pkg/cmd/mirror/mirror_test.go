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

	"github.com/stretchr/testify/require"
)

func TestModulePathToBazelRepoName(t *testing.T) {
	require.Equal(t, modulePathToBazelRepoName("github.com/alecthomas/template"), "com_github_alecthomas_template")
	require.Equal(t, modulePathToBazelRepoName("github.com/aws/aws-sdk-go-v2/service/iam"), "com_github_aws_aws_sdk_go_v2_service_iam")
	require.Equal(t, modulePathToBazelRepoName("github.com/Azure/go-ansiterm"), "com_github_azure_go_ansiterm")
	require.Equal(t, modulePathToBazelRepoName("gopkg.in/yaml.v3"), "in_gopkg_yaml_v3")
	require.Equal(t, modulePathToBazelRepoName("collectd.org"), "org_collectd")
}
