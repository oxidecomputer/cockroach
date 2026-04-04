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
	"regexp"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
)

var flowableReleaseTagRegex = regexp.MustCompile(`^flowable-(?P<major>\d+)\.(?P<minor>\d+)\.(?P<point>\d+)$`)

// This test runs Flowable test suite against a single cockroach node.

func registerFlowable(r registry.Registry) {
	runFlowable := func(
		ctx context.Context,
		t test.Test,
		c cluster.Cluster,
	) {
		if c.IsLocal() {
			t.Fatal("cannot be run in local mode")
		}
		node := c.Node(1)
		t.Status("setting up cockroach")
		c.Put(ctx, t.Cockroach(), "./cockroach", c.All())
		c.Start(ctx, t.L(), option.DefaultStartOpts(), install.MakeClusterSettings(), c.All())

		t.Status("cloning flowable and installing prerequisites")
		latestTag, err := repeatGetLatestTag(
			ctx, t, "flowable", "flowable-engine", flowableReleaseTagRegex,
		)
		if err != nil {
			t.Fatal(err)
		}
		t.L().Printf("Latest Flowable release is %s.", latestTag)

		if err := repeatRunE(
			ctx, t, c, node, "update apt-get", `sudo apt-get -qq update`,
		); err != nil {
			t.Fatal(err)
		}

		if err := repeatRunE(
			ctx,
			t,
			c,
			node,
			"install dependencies",
			`sudo apt-get -qq install default-jre openjdk-8-jdk-headless gradle maven`,
		); err != nil {
			t.Fatal(err)
		}

		if err := repeatRunE(
			ctx, t, c, node, "remove old Flowable", `rm -rf /mnt/data1/flowable-engine`,
		); err != nil {
			t.Fatal(err)
		}

		if err := repeatGitCloneE(
			ctx,
			t,
			c,
			"https://github.com/flowable/flowable-engine.git",
			"/mnt/data1/flowable-engine",
			latestTag,
			node,
		); err != nil {
			t.Fatal(err)
		}

		t.Status("building Flowable")
		if err := repeatRunE(
			ctx,
			t,
			c,
			node,
			"building Flowable",
			`cd /mnt/data1/flowable-engine/ && mvn clean install -DskipTests`,
		); err != nil {
			t.Fatal(err)
		}

		if err := c.RunE(ctx, node,
			`cd /mnt/data1/flowable-engine/ && mvn clean test -Dtest=Flowable6Test#testLongServiceTaskLoop -Ddb=crdb`,
		); err != nil {
			t.Fatal(err)
		}
	}

	r.Add(registry.TestSpec{
		Name:    "flowable",
		Owner:   registry.OwnerSQLFoundations,
		Cluster: r.MakeClusterSpec(1),
		Run: func(ctx context.Context, t test.Test, c cluster.Cluster) {
			runFlowable(ctx, t, c)
		},
	})
}
