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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/cluster"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/registry"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/spec"
	"github.com/cockroachdb/cockroach/pkg/cmd/roachtest/test"
)

func registerSyncTest(r registry.Registry) {
	const nemesisScript = `#!/usr/bin/env bash

if [[ $1 == "on" ]]; then
  charybdefs-nemesis --probability
else
  charybdefs-nemesis --clear
fi
`

	r.Add(registry.TestSpec{
		Skip:  "#48603: broken on Pebble",
		Name:  "synctest",
		Owner: registry.OwnerStorage,
		// This test sets up a custom file system; we don't want the cluster reused.
		Cluster: r.MakeClusterSpec(1, spec.ReuseNone()),
		Run: func(ctx context.Context, t test.Test, c cluster.Cluster) {
			n := c.Node(1)
			tmpDir, err := ioutil.TempDir("", "synctest")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = os.RemoveAll(tmpDir)
			}()
			nemesis := filepath.Join(tmpDir, "nemesis")

			if err := ioutil.WriteFile(nemesis, []byte(nemesisScript), 0755); err != nil {
				t.Fatal(err)
			}

			c.Put(ctx, t.Cockroach(), "./cockroach")
			c.Put(ctx, nemesis, "./nemesis")
			c.Run(ctx, n, "chmod +x nemesis")
			c.Run(ctx, n, "sudo umount {store-dir}/faulty || true")
			c.Run(ctx, n, "mkdir -p {store-dir}/{real,faulty} || true")
			t.Status("setting up charybdefs")

			if err := c.Install(ctx, t.L(), n, "charybdefs"); err != nil {
				t.Fatal(err)
			}
			c.Run(ctx, n, "sudo charybdefs {store-dir}/faulty -oallow_other,modules=subdir,subdir={store-dir}/real && chmod 777 {store-dir}/{real,faulty}")

			t.Status("running synctest")
			c.Run(ctx, n, "./cockroach debug synctest {store-dir}/faulty ./nemesis")
		},
	})
}
