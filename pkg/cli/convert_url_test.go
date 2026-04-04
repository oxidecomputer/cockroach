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

package cli

func Example_convert_url() {
	c := NewCLITest(TestCLIParams{
		NoServer: true,
	})
	defer c.Cleanup()

	c.RunWithArgs([]string{`convert-url`})
	c.RunWithArgs([]string{`convert-url`, `--url`, `postgres://foo@bar`})

	// Output:
	// convert-url
	// # WARNING: no URL specified via --url; using a random URL as example.
	//
	// # Connection URL for libpq (C/C++), psycopg (Python), lib/pq & pgx (Go), node-postgres (JS) and most pq-compatible drivers:
	// postgresql://root@localhost:26257/defaultdb
	//
	// # Connection DSN (Data Source Name) for Postgres drivers that accept DSNs - most drivers and also ODBC:
	// database=defaultdb user=root host=localhost port=26257
	//
	// # Connection URL for JDBC (Java and JVM-based languages):
	// jdbc:postgresql://localhost:26257/defaultdb?user=root
	//
	// convert-url --url postgres://foo@bar
	// # Connection URL for libpq (C/C++), psycopg (Python), lib/pq & pgx (Go), node-postgres (JS) and most pq-compatible drivers:
	// postgresql://foo@bar:26257/defaultdb
	//
	// # Connection DSN (Data Source Name) for Postgres drivers that accept DSNs - most drivers and also ODBC:
	// database=defaultdb user=foo host=bar port=26257
	//
	// # Connection URL for JDBC (Java and JVM-based languages):
	// jdbc:postgresql://bar:26257/defaultdb?user=foo
	//
}
