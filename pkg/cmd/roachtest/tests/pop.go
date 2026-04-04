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
	"context"
	"fmt"
	"regexp"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
	"github.com/stretchr/testify/require"
)

var popReleaseTag = regexp.MustCompile(`^v(?P<major>\d+)\.(?P<minor>\d+)\.(?P<point>\d+)$`)
var popSupportedTag = "v5.3.3"

func registerPop(r registry.Registry) {
	runPop := func(ctx context.Context, t test.Test, c cluster.Cluster) {
		if c.IsLocal() {
			t.Fatal("cannot be run in local mode")
		}
		node := c.Node(1)
		t.Status("setting up cockroach")
		c.Put(ctx, t.Cockroach(), "./cockroach", c.All())
		c.Start(ctx, t.L(), option.DefaultStartOpts(), install.MakeClusterSettings(), c.All())
		version, err := fetchCockroachVersion(ctx, t.L(), c, node[0])
		if err != nil {
			t.Fatal(err)
		}
		if err := alterZoneConfigAndClusterSettings(ctx, t, version, c, node[0]); err != nil {
			t.Fatal(err)
		}

		t.Status("cloning pop and installing prerequisites")
		latestTag, err := repeatGetLatestTag(
			ctx, t, "gobuffalo", "pop", popReleaseTag)
		if err != nil {
			t.Fatal(err)
		}
		t.L().Printf("Latest pop release is %s.", latestTag)
		t.L().Printf("Supported pop release is %s.", popSupportedTag)

		installGolang(ctx, t, c, node)

		const (
			popPath = "/mnt/data1/pop/"
		)

		// Remove any old pop installations
		if err := repeatRunE(
			ctx, t, c, node, "remove old pop", fmt.Sprintf("rm -rf %s", popPath),
		); err != nil {
			t.Fatal(err)
		}

		if err := repeatGitCloneE(
			ctx,
			t,
			c,
			"https://github.com/gobuffalo/pop.git",
			popPath,
			popSupportedTag,
			node,
		); err != nil {
			t.Fatal(err)
		}

		t.Status("building and setting up tests")

		err = c.RunE(ctx, node, fmt.Sprintf(`cd %s && go build -v -tags sqlite -o tsoda ./soda`, popPath))
		require.NoError(t, err)

		err = c.RunE(ctx, node, fmt.Sprintf(`cd %s && ./tsoda drop -e cockroach -c ./database.yml -p ./testdata/migrations`, popPath))
		require.NoError(t, err)

		err = c.RunE(ctx, node, fmt.Sprintf(`cd %s && ./tsoda create -e cockroach -c ./database.yml -p ./testdata/migrations`, popPath))
		require.NoError(t, err)

		err = c.RunE(ctx, node, fmt.Sprintf(`cd %s && ./tsoda migrate -e cockroach -c ./database.yml -p ./testdata/migrations`, popPath))
		require.NoError(t, err)

		t.Status("running pop test suite")

		// No tests are expected to fail.
		err = c.RunE(ctx, node, fmt.Sprintf(`cd %s && SODA_DIALECT=cockroach go test -race -tags sqlite -v ./... -count=1`, popPath))
		require.NoError(t, err, "error while running pop tests")
	}

	r.Add(registry.TestSpec{
		Name:    "pop",
		Owner:   registry.OwnerSQLFoundations,
		Cluster: r.MakeClusterSpec(1),
		Tags:    []string{`default`, `orm`},
		Run:     runPop,
	})
}
