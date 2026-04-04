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
	"math/rand"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
)

func registerRoachtest(r registry.Registry) {
	r.Add(registry.TestSpec{
		Name:    "roachtest/noop",
		Tags:    []string{"roachtest"},
		Owner:   registry.OwnerTestEng,
		Run:     func(_ context.Context, _ test.Test, _ cluster.Cluster) {},
		Cluster: r.MakeClusterSpec(0),
	})
	r.Add(registry.TestSpec{
		Name:  "roachtest/noop-maybefail",
		Tags:  []string{"roachtest"},
		Owner: registry.OwnerTestEng,
		Run: func(_ context.Context, t test.Test, _ cluster.Cluster) {
			if rand.Float64() <= 0.2 {
				t.Fatal("randomly failing")
			}
		},
		Cluster: r.MakeClusterSpec(0),
	})
	// This test can be run manually to check what happens if a test times out.
	// In particular, can manually verify that suitable artifacts are created.
	r.Add(registry.TestSpec{
		Name:  "roachtest/hang",
		Tags:  []string{"roachtest"},
		Owner: registry.OwnerTestEng,
		Run: func(_ context.Context, t test.Test, c cluster.Cluster) {
			ctx := context.Background() // intentional
			c.Put(ctx, t.Cockroach(), "cockroach", c.All())
			c.Start(ctx, t.L(), option.DefaultStartOpts(), install.MakeClusterSettings(), c.All())
			time.Sleep(time.Hour)
		},
		Timeout: 3 * time.Minute,
		Cluster: r.MakeClusterSpec(3),
	})
}
