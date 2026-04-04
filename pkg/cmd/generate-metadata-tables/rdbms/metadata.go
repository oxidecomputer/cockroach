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

// metadata.go provides Connectivity for mysql and postgres, and
// provides interfaces that helps to retrieve a schema from these
// databases for comparison purposes.

package rdbms

import (
	"context"
	"regexp"

	"github.com/cockroachdb/cockroach/pkg/sql"
)

// ConnectFns will be used to determine which kind of database will be used
// at runtime based on the flag.
var ConnectFns = map[string]func(address, user, catalog string) (DBMetadataConnection, error){
	sql.Postgres: postgresConnect,
	sql.MySQL:    mysqlConnect,
}

type columnMetadata struct {
	tableName    string
	columnName   string
	dataTypeName string
	dataTypeOid  uint32
}

type excludePattern struct {
	pattern *regexp.Regexp
	except  map[string]struct{}
}

// ColumnMetadataList is a list of rows coming from rdbms describing a column.
type ColumnMetadataList struct {
	data       []*columnMetadata
	exclusions []*excludePattern
}

// DBMetadataConnection structs can describe a schema like pg_catalog or
// information_schema.
type DBMetadataConnection interface {
	Close(ctx context.Context) error
	DescribeSchema(ctx context.Context) (*ColumnMetadataList, error)
	DatabaseVersion(ctx context.Context) (string, error)
}

// ForEachRow iterates over the rows gotten from DescribeSchema() call.
func (l *ColumnMetadataList) ForEachRow(addRow func(string, string, string, uint32)) {
	addRowIfAllowed := func(metadata *columnMetadata) {
		for _, exclusion := range l.exclusions {
			tableName := metadata.tableName
			if _, ok := exclusion.except[tableName]; exclusion.pattern.MatchString(tableName) && !ok {
				return
			}
		}

		addRow(metadata.tableName, metadata.columnName, metadata.dataTypeName, metadata.dataTypeOid)
	}
	for _, metadata := range l.data {
		addRowIfAllowed(metadata)
	}
}
