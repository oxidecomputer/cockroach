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

package sql_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/tests"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// Use the output from crdb_internal.show_create_all_tables() to recreate the
// tables and perform another crdb_internal.show_create_all_tables() to ensure
// that the output is the same after recreating the tables.
func TestRecreateTables(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	params, _ := tests.CreateTestServerParams()
	s, sqlDB, _ := serverutils.StartServer(t, params)
	defer s.Stopper().Stop(context.Background())

	sqlRunner := sqlutils.MakeSQLRunner(sqlDB)
	sqlRunner.Exec(t, `CREATE DATABASE test;`)
	sqlRunner.Exec(t, `USE test;`)
	sqlRunner.Exec(t, `CREATE TABLE foo(x INT primary key);`)
	sqlRunner.Exec(t, `CREATE TABLE bar(x INT, y INT, z STRING, FAMILY f1(x, y, z))`)

	row := sqlRunner.QueryRow(t, "SELECT crdb_internal.show_create_all_tables('test')")
	var recreateTablesStmt string
	row.Scan(&recreateTablesStmt)

	// Use the recreateTablesStmt to recreate the tables, perform another
	// show_create_all_tables and compare that the output is the same.
	sqlRunner.Exec(t, `DROP DATABASE test;`)
	sqlRunner.Exec(t, `CREATE DATABASE test;`)
	sqlRunner.Exec(t, recreateTablesStmt)

	row = sqlRunner.QueryRow(t, "SELECT crdb_internal.show_create_all_tables('test')")
	var recreateTablesStmt2 string
	row.Scan(&recreateTablesStmt2)

	if recreateTablesStmt != recreateTablesStmt2 {
		t.Fatalf("got: %s\nexpected: %s", recreateTablesStmt2, recreateTablesStmt)
	}
}
