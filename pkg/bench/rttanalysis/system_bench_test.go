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

package rttanalysis

import "testing"

func BenchmarkSystemDatabaseQueries(b *testing.B) { reg.Run(b) }
func init() {
	reg.Register("SystemDatabaseQueries", []RoundTripBenchTestCase{
		// This query performs 1-2 lookups: getting the descriptor ID by Name, then
		// fetching the system table descriptor. The descriptor is then cached.
		{
			Name: "select system.users with schema Name",
			Stmt: `SELECT username, "hashedPassword" FROM system.public.users WHERE username = 'root'`,
		},
		// This query performs 1 extra lookup since the executor first tries to
		// lookup the Name `current_db.system.users`.
		{
			Name: "select system.users without schema Name",
			Stmt: `SELECT username, "hashedPassword" FROM system.users WHERE username = 'root'`,
		},
		// This query performs 0 extra lookups since the Name resolution logic does
		// not try to resolve `"".system.users` and instead resolves
		//`system.public.users` right away.
		{
			Name:  "select system.users with empty database Name",
			Setup: `SET sql_safe_updates = false; USE "";`,
			Stmt:  `SELECT username, "hashedPassword"  FROM system.users WHERE username = 'root'`,
		},
	})
}
