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

func BenchmarkVirtualTableQueries(b *testing.B) { reg.Run(b) }
func init() {
	reg.Register("VirtualTableQueries", []RoundTripBenchTestCase{
		// This benchmark should perform exactly one kv operation to fetch all of
		// the descriptors.
		{
			Name: "select crdb_internal.tables with 1 fk",
			Setup: `
CREATE TABLE t1 (i INT PRIMARY KEY);
CREATE TABLE t2 (i INT PRIMARY KEY, j INT REFERENCES t1(i));
`,
			Stmt: `SELECT * FROM "".crdb_internal.tables`,
		},
		{
			// We expect this to perform exactly one kv operation to fetch all of
			// the descriptors.
			Name: "select crdb_internal.invalid_objects with 1 fk",
			Setup: `
CREATE TABLE t1 (i INT PRIMARY KEY);
CREATE TABLE t2 (i INT PRIMARY KEY, j INT REFERENCES t1(i));
`,
			Stmt: `SELECT * FROM "".crdb_internal.invalid_objects`,
		},
	})
}
