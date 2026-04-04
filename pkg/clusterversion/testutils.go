// Copyright 2019 The Cockroach Authors.
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

// TestingBinaryVersion is a binary version that tests can use when they don't
// want to go through a Settings object.
var TestingBinaryVersion = binaryVersion

// TestingBinaryMinSupportedVersion is a minimum supported version that
// tests can use when they don't want to go through a Settings object.
var TestingBinaryMinSupportedVersion = binaryMinSupportedVersion

// TestingClusterVersion is a ClusterVersion that tests can use when they don't
// want to go through a Settings object.
var TestingClusterVersion = ClusterVersion{
	Version: TestingBinaryVersion,
}
