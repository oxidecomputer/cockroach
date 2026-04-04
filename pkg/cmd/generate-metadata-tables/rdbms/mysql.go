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

// mysql.go has the implementations of DBMetadataConnection to
// connect and retrieve schemas from mysql rdbms.

package rdbms

import (
	"context"
	gosql "database/sql"
	"fmt"
	"regexp"
	"strings"

	// gosql implementation.
	_ "github.com/go-sql-driver/mysql"
)

const mysqlDescribeSchema = `
	SELECT 
		table_name, 
		column_name, 
		data_type 
	FROM information_schema.columns
	WHERE table_schema = ?
	ORDER BY table_name
`

var mysqlExclusions = []*excludePattern{
	{
		pattern: regexp.MustCompile(`innodb_.+`),
		except:  make(map[string]struct{}),
	},
}

type mysqlMetadataConnection struct {
	*gosql.DB
	catalog string
}

func mysqlConnect(address, user, catalog string) (DBMetadataConnection, error) {
	db, err := gosql.Open("mysql", fmt.Sprintf("%s@tcp(%s)/%s", user, address, catalog))
	if err != nil {
		return nil, err
	}
	return mysqlMetadataConnection{db, catalog}, nil
}

func (conn mysqlMetadataConnection) DatabaseVersion(
	ctx context.Context,
) (version string, err error) {
	row := conn.QueryRowContext(ctx, "SELECT version()")
	err = row.Scan(&version)
	return version, err
}

func (conn mysqlMetadataConnection) DescribeSchema(
	ctx context.Context,
) (*ColumnMetadataList, error) {
	metadata := &ColumnMetadataList{exclusions: mysqlExclusions}
	rows, err := conn.QueryContext(ctx, mysqlDescribeSchema, conn.catalog)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row := new(columnMetadata)
		if err := rows.Scan(&row.tableName, &row.columnName, &row.dataTypeName); err != nil {
			return nil, err
		}
		row.tableName = strings.ToLower(row.tableName)
		row.columnName = strings.ToLower(row.columnName)
		metadata.data = append(metadata.data, row)
	}

	return metadata, nil
}

func (conn mysqlMetadataConnection) Close(ctx context.Context) error {
	return conn.DB.Close()
}
