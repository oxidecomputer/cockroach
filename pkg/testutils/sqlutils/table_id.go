// Copyright 2016 The Cockroach Authors.
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

package sqlutils

import (
	"context"
	"testing"
)

// QueryDatabaseID returns the database ID of the specified database using the
// system.namespace table.
func QueryDatabaseID(t testing.TB, sqlDB DBHandle, dbName string) uint32 {
	dbIDQuery := `
		SELECT id FROM system.namespace
		WHERE name = $1 AND "parentSchemaID" = 0 AND "parentID" = 0
	`
	var dbID uint32
	result := sqlDB.QueryRowContext(context.Background(), dbIDQuery, dbName)
	if err := result.Scan(&dbID); err != nil {
		t.Fatal(err)
	}
	return dbID
}

// QuerySchemaID returns the schema ID of the specified database.schema
// using the system.namespace table.
func QuerySchemaID(t testing.TB, sqlDB DBHandle, dbName, schemaName string) uint32 {
	tableIDQuery := `
 SELECT schemas.id FROM system.namespace schemas
   JOIN system.namespace dbs ON dbs.id = schemas."parentID"
   WHERE dbs.name = $1 AND schemas.name = $2
 `
	var schemaID uint32
	result := sqlDB.QueryRowContext(
		context.Background(),
		tableIDQuery, dbName,
		schemaName,
	)
	if err := result.Scan(&schemaID); err != nil {
		t.Fatal(err)
	}
	return schemaID
}

// QueryTableID returns the table ID of the specified database.table
// using the system.namespace table.
func QueryTableID(
	t testing.TB, sqlDB DBHandle, dbName, schemaName string, tableName string,
) uint32 {
	tableIDQuery := `
 SELECT tables.id FROM system.namespace tables
   JOIN system.namespace dbs ON dbs.id = tables."parentID"
	 JOIN system.namespace schemas ON schemas.id = tables."parentSchemaID"
   WHERE dbs.name = $1 AND schemas.name = $2 AND tables.name = $3
 `
	var tableID uint32
	result := sqlDB.QueryRowContext(
		context.Background(),
		tableIDQuery, dbName,
		schemaName,
		tableName,
	)
	if err := result.Scan(&tableID); err != nil {
		t.Fatal(err)
	}
	return tableID
}
