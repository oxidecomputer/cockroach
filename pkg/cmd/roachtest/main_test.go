// Copyright 2022 The Cockroach Authors.
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

import "github.com/cockroachdb/cockroach/pkg/internal/team"

func init() {
	// Set a bogus build tag. Tests that make a testRegistry would otherwise end
	// up shelling out to `git` which may not work (for example if the tests are
	// run through bazel).
	buildTag = "v99.99.99"
	// Similar for TEAMS.yaml.
	loadTeams = func() (team.Map, error) {
		return map[team.Alias]team.Team{
			ownerToAlias(OwnerUnitTest): {},
		}, nil
	}
}
