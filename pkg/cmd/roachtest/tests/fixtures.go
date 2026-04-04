// Copyright 2018 The Cockroach Authors.
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
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
)

func registerFixtures(r registry.Registry) {
	// Run this test to create a new fixture for the version upgrade test. This
	// is necessary after every release. For example, the day `master` becomes
	// the 20.2 release, this test will fail because it is missing a fixture for
	// 20.1; run the test (on 20.1). Check it in (instructions will be logged
	// below) and off we go.
	//
	// The version to create/update the fixture for. Must be released (i.e.
	// can download it from the homepage); if that is not the case use the
	// empty string which uses the local cockroach binary. Make sure that
	// this binary then has the correct version. For example, to make a
	// "v20.2" fixture, you will need a binary that has "v20.2" in the
	// output of `./cockroach version`, and this process will end up
	// creating fixtures that have "v20.2" in them. This would be part
	// of tagging the master branch as v21.1 in the process of going
	// through the major release for v20.2. The version is passed in as
	// FIXTURE_VERSION environment variable.
	//
	// In the common case, one should populate this with the version (instead of
	// using the empty string) as this is the most straightforward and least
	// error-prone way to generate the fixtures.
	//
	// Please note that you do *NOT* need to update the fixtures in a patch
	// release. This only happens as part of preparing the master branch for the
	// next release. The release team runbooks, at time of writing, reflect
	// this.
	//
	// Example invocation:
	// roachtest --local run generate-fixtures --debug --cockroach ./cockroach \
	//   --build-tag v22.1.0-beta.3 tag:fixtures
	runFixtures := func(
		ctx context.Context,
		t test.Test,
		c cluster.Cluster,
	) {
		fixtureVersion := strings.TrimPrefix(t.BuildVersion().String(), "v")
		makeVersionFixtureAndFatal(ctx, t, c, fixtureVersion)
	}
	spec := registry.TestSpec{
		Name:    "generate-fixtures",
		Timeout: 30 * time.Minute,
		Tags:    []string{"fixtures"},
		Owner:   registry.OwnerDevInf,
		Cluster: r.MakeClusterSpec(4),
		Run:     runFixtures,
	}
	r.Add(spec)
}
