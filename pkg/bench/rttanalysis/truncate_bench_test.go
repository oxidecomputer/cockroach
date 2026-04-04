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

func BenchmarkTruncate(b *testing.B) { reg.Run(b) }
func init() {
	reg.Register("Truncate", []RoundTripBenchTestCase{
		{
			Name:  "truncate 1 column 0 rows",
			Setup: "CREATE TABLE t(x INT);",
			Stmt:  "TRUNCATE t",
		},
		{
			Name: "truncate 1 column 1 row",
			Setup: `CREATE TABLE t(x INT); 
INSERT INTO t (x) VALUES (1);`,
			Stmt: "TRUNCATE t",
		},
		{
			Name: "truncate 1 column 2 rows",
			Setup: `CREATE TABLE t(x INT); 
INSERT INTO t (x) VALUES (1);
INSERT INTO t (x) VALUES (2);`,
			Stmt: "TRUNCATE t",
		},
		{
			Name:  "truncate 2 column 0 rows",
			Setup: `CREATE TABLE t(x INT, y INT);`,
			Stmt:  "TRUNCATE t",
		},
		{
			Name: "truncate 2 column 1 rows",
			Setup: `CREATE TABLE t(x INT, y INT); 
INSERT INTO t (x, y) VALUES (1, 1);`,
			Stmt: "TRUNCATE t",
		},
		{
			Name: "truncate 2 column 2 rows",
			Setup: `CREATE TABLE t(x INT, y INT); 
INSERT INTO t (x, y) VALUES (1, 1);
INSERT INTO t (x,y) VALUES (2, 2);`,
			Stmt: "TRUNCATE t",
		},
	})
}
