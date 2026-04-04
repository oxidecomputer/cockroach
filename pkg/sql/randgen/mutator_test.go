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

package randgen

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
)

func TestPostgresMutator(t *testing.T) {
	defer leaktest.AfterTest(t)()
	q := `
		CREATE TABLE t (s STRING FAMILY fam1, b BYTES, FAMILY fam2 (b), PRIMARY KEY (s ASC, b DESC), INDEX (s) STORING (b))
		    PARTITION BY LIST (s)
		        (
		            PARTITION europe_west VALUES IN ('a', 'b')
		        );
		ALTER TABLE table1 INJECT STATISTICS 'blah';
		SET CLUSTER SETTING "sql.stats.automatic_collection.enabled" = false;
	`

	rng, _ := randutil.NewTestRand()
	{
		mutated, changed := ApplyString(rng, q, PostgresMutator)
		if !changed {
			t.Fatal("expected changed")
		}
		mutated = strings.TrimSpace(mutated)
		expect := `CREATE TABLE t (s TEXT, b BYTEA, PRIMARY KEY (s ASC, b DESC), INDEX (s) INCLUDE (b));`
		if mutated != expect {
			t.Fatalf("unexpected: %s", mutated)
		}
	}
	{
		mutated, changed := ApplyString(rng, q, PostgresCreateTableMutator, PostgresMutator)
		if !changed {
			t.Fatal("expected changed")
		}
		mutated = strings.TrimSpace(mutated)
		expect := "CREATE TABLE t (s TEXT, b BYTEA, PRIMARY KEY (s, b));\nCREATE INDEX ON t (s) INCLUDE (b);"
		if mutated != expect {
			t.Fatalf("unexpected: %s", mutated)
		}
	}
}
