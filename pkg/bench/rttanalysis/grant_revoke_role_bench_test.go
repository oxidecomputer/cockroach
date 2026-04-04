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

func BenchmarkGrantRole(b *testing.B) { reg.Run(b) }
func init() {
	reg.Register("GrantRole", []RoundTripBenchTestCase{
		{
			Name: "grant 1 role",
			Setup: `CREATE ROLE a;
CREATE ROLE b;`,
			Stmt:  "GRANT a TO b",
			Reset: "DROP ROLE a,b",
		},
		{
			Name: "grant 2 roles",
			Setup: `CREATE ROLE a;
CREATE ROLE b;
CREATE ROLE c;`,
			Stmt:  "GRANT a,b TO c",
			Reset: "DROP ROLE a,b,c",
		},
	})
}

func BenchmarkRevokeRole(b *testing.B) { reg.Run(b) }
func init() {
	reg.Register("RevokeRole", []RoundTripBenchTestCase{
		{
			Name: "revoke 1 role",
			Setup: `CREATE ROLE a;
CREATE ROLE b;
GRANT a TO b`,
			Stmt:  "REVOKE a FROM b",
			Reset: "DROP ROLE a,b",
		},
		{
			Name: "revoke 2 roles",
			Setup: `CREATE ROLE a;
CREATE ROLE b;
CREATE ROLE c;
GRANT a,b TO c;`,
			Stmt:  "REVOKE a,b FROM c",
			Reset: "DROP ROLE a,b,c",
		},
	})
}
