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

package workloadsql

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
	"github.com/cockroachdb/cockroach/pkg/workload"
	"github.com/cockroachdb/cockroach/pkg/workload/bank"
)

func TestSetup(t *testing.T) {
	defer leaktest.AfterTest(t)()

	tests := []struct {
		rows        int
		batchSize   int
		concurrency int
	}{
		{10, 1, 1},
		{10, 9, 1},
		{10, 10, 1},
		{10, 100, 1},
		{10, 1, 4},
		{10, 9, 4},
		{10, 10, 4},
		{10, 100, 4},
	}

	ctx := context.Background()
	s, db, _ := serverutils.StartServer(t, base.TestServerArgs{UseDatabase: `test`})
	defer s.Stopper().Stop(ctx)
	sqlutils.MakeSQLRunner(db).Exec(t, `CREATE DATABASE test`)

	for _, test := range tests {
		t.Run(fmt.Sprintf("rows=%d/batch=%d", test.rows, test.batchSize), func(t *testing.T) {
			sqlDB := sqlutils.MakeSQLRunner(db)
			sqlDB.Exec(t, `DROP TABLE IF EXISTS bank`)

			gen := bank.FromRows(test.rows)
			l := InsertsDataLoader{BatchSize: test.batchSize, Concurrency: test.concurrency}
			if _, err := Setup(ctx, db, gen, l); err != nil {
				t.Fatalf("%+v", err)
			}

			for _, table := range gen.Tables() {
				var c int
				sqlDB.QueryRow(t, fmt.Sprintf(`SELECT count(*) FROM %s`, table.Name)).Scan(&c)
				// There happens to be 1 row per batch in bank.
				if c != table.InitialRows.NumBatches {
					t.Errorf(`%s: got %d rows expected %d`,
						table.Name, c, table.InitialRows.NumBatches)
				}
			}
		})
	}
}

func TestSplits(t *testing.T) {
	defer leaktest.AfterTest(t)()

	ctx := context.Background()
	s, db, _ := serverutils.StartServer(t, base.TestServerArgs{UseDatabase: `test`})
	defer s.Stopper().Stop(ctx)
	sqlutils.MakeSQLRunner(db).Exec(t, `CREATE DATABASE test`)

	for _, ranges := range []int{1, 2, 3, 4, 10} {

		tables := []workload.Table{
			{
				Name:   `ints`,
				Schema: `(a INT PRIMARY KEY)`,
				Splits: workload.Tuples(ranges-1, func(i int) []interface{} {
					return []interface{}{i}
				}),
			},
			{
				Name:   `floats`,
				Schema: `(a FLOAT PRIMARY KEY)`,
				Splits: workload.Tuples(ranges-1, func(i int) []interface{} {
					return []interface{}{float64(i)}
				}),
			},
			{
				Name:   `strings`,
				Schema: `(a STRING PRIMARY KEY)`,
				Splits: workload.Tuples(ranges-1, func(i int) []interface{} {
					return []interface{}{strconv.Itoa(i)}
				}),
			},
			{
				Name:   `bytes`,
				Schema: `(a BYTES PRIMARY KEY)`,
				Splits: workload.Tuples(ranges-1, func(i int) []interface{} {
					return []interface{}{strconv.Itoa(i)}
				}),
			},
			{
				Name:   `uuids`,
				Schema: `(a UUID PRIMARY KEY)`,
				Splits: workload.Tuples(ranges-1, func(i int) []interface{} {
					u, err := uuid.NewV4()
					if err != nil {
						panic(err)
					}
					return []interface{}{u.String()}
				}),
			},
		}

		t.Run(fmt.Sprintf("ranges=%d", ranges), func(t *testing.T) {
			sqlDB := sqlutils.MakeSQLRunner(db)
			for _, table := range tables {
				sqlDB.Exec(t, fmt.Sprintf(`DROP TABLE IF EXISTS %s`, table.Name))
				sqlDB.Exec(t, fmt.Sprintf(`CREATE TABLE %s %s`, table.Name, table.Schema))

				const concurrency = 10
				if err := Split(ctx, db, table, concurrency); err != nil {
					t.Fatalf("%+v", err)
				}

				countRangesQ := fmt.Sprintf(
					`SELECT count(*) FROM [SHOW RANGES FROM TABLE test.%s]`, table.Name,
				)
				var actual int
				sqlDB.QueryRow(t, countRangesQ).Scan(&actual)
				if ranges != actual {
					t.Errorf(`expected %d got %d`, ranges, actual)
				}
			}
		})
	}
}
