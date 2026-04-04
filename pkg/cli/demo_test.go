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

package cli

import (
	"github.com/cockroachdb/cockroach/pkg/cli/democluster"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func Example_demo() {
	c := NewCLITest(TestCLIParams{NoServer: true})
	defer c.Cleanup()

	defer democluster.TestingForceRandomizeDemoPorts()()

	// Disable multi-tenant, as it requires a CCL binary
	testData := [][]string{
		{`demo`, `--multitenant=false`, `-e`, `show database`},
		{`demo`, `--multitenant=false`, `-e`, `show database`, `--no-example-database`},
		{`demo`, `--multitenant=false`, `-e`, `show application_name`},
		{`demo`, `--multitenant=false`, `--format=table`, `-e`, `show database`},
		{`demo`, `--multitenant=false`, `-e`, `select 1 as "1"`, `-e`, `select 3 as "3"`},
		{`demo`, `--multitenant=false`, `--echo-sql`, `-e`, `select 1 as "1"`},
		{`demo`, `--multitenant=false`, `--set=errexit=0`, `-e`, `select nonexistent`, `-e`, `select 123 as "123"`},
		{`demo`, `--multitenant=false`, `startrek`, `-e`, `SELECT database_name, owner FROM [show databases]`},
		{`demo`, `--multitenant=false`, `startrek`, `-e`, `SELECT database_name, owner FROM [show databases]`, `--format=table`},
		// Test that if we start with --insecure we cannot perform
		// commands that require a secure cluster.
		{`demo`, `--multitenant=false`, `-e`, `CREATE USER test WITH PASSWORD 'testpass'`},
		{`demo`, `--multitenant=false`, `--insecure`, `-e`, `CREATE USER test WITH PASSWORD 'testpass'`},
		{`demo`, `--multitenant=false`, `--geo-partitioned-replicas`, `--disable-demo-license`},
	}
	setCLIDefaultsForTests()
	// We must reset the security asset loader here, otherwise the dummy
	// asset loader that is set by default in tests will not be able to
	// find the certs that demo sets up.
	security.ResetAssetLoader()
	for _, cmd := range testData {
		// `demo` sets up a server and log file redirection, which asserts
		// that the logging subsystem has not been initialized yet. Fake
		// this to be true.
		log.TestingResetActive()
		c.RunWithArgs(cmd)
	}

	// Output:
	// demo --multitenant=false -e show database
	// database
	// movr
	// demo --multitenant=false -e show database --no-example-database
	// database
	// defaultdb
	// demo --multitenant=false -e show application_name
	// application_name
	// $ cockroach demo
	// demo --multitenant=false --format=table -e show database
	//   database
	// ------------
	//   movr
	// (1 row)
	// demo --multitenant=false -e select 1 as "1" -e select 3 as "3"
	// 1
	// 1
	// 3
	// 3
	// demo --multitenant=false --echo-sql -e select 1 as "1"
	// > select 1 as "1"
	// 1
	// 1
	// demo --multitenant=false --set=errexit=0 -e select nonexistent -e select 123 as "123"
	// ERROR: column "nonexistent" does not exist
	// SQLSTATE: 42703
	// 123
	// 123
	// demo --multitenant=false startrek -e SELECT database_name, owner FROM [show databases]
	// database_name	owner
	// defaultdb	root
	// postgres	root
	// startrek	demo
	// system	node
	// demo --multitenant=false startrek -e SELECT database_name, owner FROM [show databases] --format=table
	//   database_name | owner
	// ----------------+--------
	//   defaultdb     | root
	//   postgres      | root
	//   startrek      | demo
	//   system        | node
	// (4 rows)
	// demo --multitenant=false -e CREATE USER test WITH PASSWORD 'testpass'
	// CREATE ROLE
	// demo --multitenant=false --insecure -e CREATE USER test WITH PASSWORD 'testpass'
	// ERROR: setting or updating a password is not supported in insecure mode
	// SQLSTATE: 28P01
	// demo --multitenant=false --geo-partitioned-replicas --disable-demo-license
	// ERROR: enterprise features are needed for this demo (--geo-partitioned-replicas)
}
