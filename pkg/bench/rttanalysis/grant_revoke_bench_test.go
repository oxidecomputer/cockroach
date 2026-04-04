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

package rttanalysis

import "testing"

func BenchmarkGrant(b *testing.B) { reg.Run(b) }
func init() {
	reg.Register("Grant", []RoundTripBenchTestCase{
		{
			Name: "grant all on 1 table",
			Setup: `CREATE USER TEST; 
CREATE TABLE t0();`,
			Stmt:  "GRANT ALL ON * TO TEST",
			Reset: "DROP ROLE TEST",
		},
		{
			Name: "grant all on 2 tables",
			Setup: `CREATE USER TEST; 
CREATE TABLE t0();
CREATE TABLE t1();`,
			Stmt:  "GRANT ALL ON * TO TEST",
			Reset: "DROP ROLE TEST",
		},
		{
			Name: "grant all on 3 tables",
			Setup: `CREATE USER TEST; 
CREATE TABLE t0();
CREATE TABLE t1();
CREATE TABLE t2();`,
			Stmt:  "GRANT ALL ON * TO TEST",
			Reset: "DROP ROLE TEST",
		},
	})
}

func BenchmarkRevoke(b *testing.B) { reg.Run(b) }
func init() {
	reg.Register("Revoke", []RoundTripBenchTestCase{
		{
			Name: "revoke all on 1 table",
			Setup: `CREATE USER TEST; 
CREATE TABLE t0();
GRANT ALL ON * TO TEST;`,
			Stmt:  "REVOKE ALL ON * FROM TEST",
			Reset: "DROP ROLE TEST",
		},
		{
			Name: "revoke all on 2 tables",
			Setup: `CREATE USER TEST; 
CREATE TABLE t0();
CREATE TABLE t1();
GRANT ALL ON * TO TEST;`,
			Stmt:  "REVOKE ALL ON * FROM TEST",
			Reset: "DROP ROLE TEST",
		},
		{
			Name: "revoke all on 3 tables",
			Setup: `CREATE USER TEST; 
CREATE TABLE t0();
CREATE TABLE t1();
CREATE TABLE t2();
GRANT ALL ON * TO TEST;`,
			Stmt:  "REVOKE ALL ON * FROM TEST",
			Reset: "DROP ROLE TEST",
		},
	})
}
