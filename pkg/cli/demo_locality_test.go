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

//go:build !race
// +build !race

package cli

import (
	"github.com/cockroachdb/cockroach/pkg/cli/democluster"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func Example_demo_locality() {
	c := NewCLITest(TestCLIParams{NoServer: true})
	defer c.Cleanup()

	defer democluster.TestingForceRandomizeDemoPorts()()

	// Disable multi-tenant for this test due to the unsupported gossip commands.
	testData := [][]string{
		{`demo`, `--nodes`, `3`, `--multitenant=false`, `-e`, `select node_id, locality from crdb_internal.gossip_nodes order by node_id`},
		{`demo`, `--nodes`, `9`, `--multitenant=false`, `-e`, `select node_id, locality from crdb_internal.gossip_nodes order by node_id`},
		{`demo`, `--nodes`, `3`, `--multitenant=false`, `--demo-locality=region=us-east1:region=us-east2:region=us-east3`,
			`-e`, `select node_id, locality from crdb_internal.gossip_nodes order by node_id`},
	}
	setCLIDefaultsForTests()
	// We must reset the security asset loader here, otherwise the dummy
	// asset loader that is set by default in tests will not be able to
	// find the certs that demo sets up.
	security.ResetAssetLoader()
	for _, cmd := range testData {
		// `demo` sets up a server and log file redirection, which asserts
		// that the logging subsystem has not been initialized yet.  Fake
		// this to be true.
		log.TestingResetActive()
		c.RunWithArgs(cmd)
	}

	// Output:
	// demo --nodes 3 --multitenant=false -e select node_id, locality from crdb_internal.gossip_nodes order by node_id
	// node_id	locality
	// 1	region=us-east1,az=b
	// 2	region=us-east1,az=c
	// 3	region=us-east1,az=d
	// demo --nodes 9 --multitenant=false -e select node_id, locality from crdb_internal.gossip_nodes order by node_id
	// node_id	locality
	// 1	region=us-east1,az=b
	// 2	region=us-east1,az=c
	// 3	region=us-east1,az=d
	// 4	region=us-west1,az=a
	// 5	region=us-west1,az=b
	// 6	region=us-west1,az=c
	// 7	region=europe-west1,az=b
	// 8	region=europe-west1,az=c
	// 9	region=europe-west1,az=d
	// demo --nodes 3 --multitenant=false --demo-locality=region=us-east1:region=us-east2:region=us-east3 -e select node_id, locality from crdb_internal.gossip_nodes order by node_id
	// node_id	locality
	// 1	region=us-east1
	// 2	region=us-east2
	// 3	region=us-east3
}
