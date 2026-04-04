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
	"fmt"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/option"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
	"github.com/cockroachdb/cockroach/pkg/roachprod/install"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/cockroach/pkg/util/version"
	"github.com/cockroachdb/errors"
	"golang.org/x/exp/rand"
)

func registerEngineSwitch(r registry.Registry) {
	runEngineSwitch := func(ctx context.Context, t test.Test, c cluster.Cluster, additionalArgs ...string) {
		roachNodes := c.Range(1, c.Spec().NodeCount-1)
		loadNode := c.Node(c.Spec().NodeCount)
		c.Put(ctx, t.DeprecatedWorkload(), "./workload", loadNode)
		c.Put(ctx, t.Cockroach(), "./cockroach", roachNodes)

		rockdbStartOpts := option.DefaultStartOpts()
		rockdbStartOpts.RoachprodOpts.ExtraArgs = append(rockdbStartOpts.RoachprodOpts.ExtraArgs, "--storage-engine=rocksdb")

		pebbleStartOpts := option.DefaultStartOpts()
		pebbleStartOpts.RoachprodOpts.ExtraArgs = append(pebbleStartOpts.RoachprodOpts.ExtraArgs, "--storage-engine=pebble")
		c.Start(ctx, t.L(), rockdbStartOpts, install.MakeClusterSettings(), roachNodes)
		stageDuration := 1 * time.Minute
		if c.IsLocal() {
			t.L().Printf("local mode: speeding up test\n")
			stageDuration = 10 * time.Second
		}
		numIters := 5 * len(roachNodes)

		loadDuration := " --duration=" + (time.Duration(numIters) * stageDuration).String()

		var deprecatedWorkloadsStr string
		if !t.BuildVersion().AtLeast(version.MustParse("v20.2.0")) {
			deprecatedWorkloadsStr += " --deprecated-fk-indexes"
		}

		workloads := []string{
			// Currently tpcc is the only one with CheckConsistency. We can add more later.
			"./workload run tpcc --tolerate-errors --wait=false --drop --init" + deprecatedWorkloadsStr + " --warehouses=1 " + loadDuration + " {pgurl:1-%d}",
		}
		checkWorkloads := []string{
			"./workload check tpcc --warehouses=1 --expensive-checks=true {pgurl:1}",
		}
		m := c.NewMonitor(ctx, roachNodes)
		for _, cmd := range workloads {
			cmd := cmd // loop-local copy
			m.Go(func(ctx context.Context) error {
				cmd = fmt.Sprintf(cmd, len(roachNodes))
				return c.RunE(ctx, loadNode, cmd)
			})
		}

		usingPebble := make([]bool, len(roachNodes))
		rng := rand.New(rand.NewSource(uint64(timeutil.Now().UnixNano())))
		m.Go(func(ctx context.Context) error {
			l, err := t.L().ChildLogger("engine-switcher")
			if err != nil {
				return err
			}
			// NB: the number of calls to `sleep` needs to be reflected in `loadDuration`.
			sleepAndCheck := func() error {
				t.WorkerStatus("sleeping")
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(stageDuration):
				}
				// Make sure everyone is still running.
				for i := 1; i <= len(roachNodes); i++ {
					t.WorkerStatus("checking ", i)
					db := c.Conn(ctx, t.L(), i)
					defer db.Close()
					rows, err := db.Query(`SHOW DATABASES`)
					if err != nil {
						return err
					}
					if err := rows.Close(); err != nil {
						return err
					}
					if err := c.CheckReplicaDivergenceOnDB(ctx, t.L(), db); err != nil {
						return errors.Wrapf(err, "node %d", i)
					}
				}
				return nil
			}

			for i := 0; i < numIters; i++ {
				// First let the load generators run in the cluster.
				if err := sleepAndCheck(); err != nil {
					return err
				}

				stop := func(node int) error {
					m.ExpectDeath()
					if rng.Intn(2) == 0 {
						l.Printf("stopping node gracefully %d\n", node)
						return c.StopCockroachGracefullyOnNode(ctx, t.L(), node)
					}
					l.Printf("stopping node %d\n", node)
					c.Stop(ctx, t.L(), option.DefaultStopOpts(), c.Node(node))
					return nil
				}

				i := rng.Intn(len(roachNodes))
				var opts option.StartOpts
				usingPebble[i] = !usingPebble[i]
				if usingPebble[i] {
					opts = pebbleStartOpts
				} else {
					opts = rockdbStartOpts
				}
				t.WorkerStatus("switching ", i+1)
				l.Printf("switching %d\n", i+1)
				if err := stop(i + 1); err != nil {
					return err
				}
				c.Start(ctx, t.L(), opts, install.MakeClusterSettings(), c.Node(i+1))
			}
			return sleepAndCheck()
		})
		m.Wait()

		for _, cmd := range checkWorkloads {
			c.Run(ctx, loadNode, cmd)
		}
	}

	n := 3
	r.Add(registry.TestSpec{
		Name:    fmt.Sprintf("engine/switch/nodes=%d", n),
		Owner:   registry.OwnerStorage,
		Skip:    "rocksdb removed in 21.1",
		Cluster: r.MakeClusterSpec(n + 1),
		Run: func(ctx context.Context, t test.Test, c cluster.Cluster) {
			runEngineSwitch(ctx, t, c)
		},
	})
	r.Add(registry.TestSpec{
		Name:    fmt.Sprintf("engine/switch/encrypted/nodes=%d", n),
		Owner:   registry.OwnerStorage,
		Skip:    "rocksdb removed in 21.1",
		Cluster: r.MakeClusterSpec(n + 1),
		Run: func(ctx context.Context, t test.Test, c cluster.Cluster) {
			runEngineSwitch(ctx, t, c, "--encrypt=true")
		},
	})
}
