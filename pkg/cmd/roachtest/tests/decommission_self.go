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

package tests

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
)

// runDecommissionSelf decommissions n2 through n2. This is an acceptance test.
//
// See https://github.com/cockroachdb/cockroach/issues/56718
func runDecommissionSelf(ctx context.Context, t test.Test, c cluster.Cluster) {
	// An empty string means that the cockroach binary specified by flag
	// `cockroach` will be used.
	const mainVersion = ""

	allNodes := c.All()
	u := newVersionUpgradeTest(c,
		uploadVersionStep(allNodes, mainVersion),
		startVersion(allNodes, mainVersion),
		fullyDecommissionStep(2, 2, mainVersion),
		func(ctx context.Context, t test.Test, u *versionUpgradeTest) {
			// Stop n2 and exclude it from post-test consistency checks,
			// as this node can't contact cluster any more and operations
			// on it will hang.
			u.c.Wipe(ctx, c.Node(2))
		},
		checkOneMembership(1, "decommissioned"),
	)

	u.run(ctx, t)
}
