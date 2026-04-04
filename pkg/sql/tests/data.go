// Copyright 2017 The Cockroach Authors.
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

package tests

import (
	"bytes"
	"context"
	gosql "database/sql"
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/errors"
)

// CheckKeyCount checks that the number of keys in the provided span matches
// numKeys.
func CheckKeyCount(t *testing.T, kvDB *kv.DB, span roachpb.Span, numKeys int) {
	t.Helper()
	if err := CheckKeyCountE(t, kvDB, span, numKeys); err != nil {
		t.Fatal(err)
	}
}

// CheckKeyCountE returns an error if the the number of keys in the
// provided span does not match numKeys.
func CheckKeyCountE(t *testing.T, kvDB *kv.DB, span roachpb.Span, numKeys int) error {
	t.Helper()
	if kvs, err := kvDB.Scan(context.TODO(), span.Key, span.EndKey, 0); err != nil {
		return err
	} else if l := numKeys; len(kvs) != l {
		return errors.Newf("expected %d key value pairs, but got %d", l, len(kvs))
	}
	return nil
}

// CreateKVTable creates a basic table named t.<name> that stores key/value
// pairs with numRows of arbitrary data.
func CreateKVTable(sqlDB *gosql.DB, name string, numRows int) error {
	// Fix the column families so the key counts don't change if the family
	// heuristics are updated.
	schemaStmts := []string{
		`CREATE DATABASE IF NOT EXISTS t;`,
		fmt.Sprintf(`CREATE TABLE t.%s (k INT PRIMARY KEY, v INT, FAMILY (k), FAMILY (v));`, name),
		fmt.Sprintf(`CREATE INDEX foo on t.%s (v);`, name),
	}

	for _, stmt := range schemaStmts {
		if _, err := sqlDB.Exec(stmt); err != nil {
			return err
		}
	}

	// Bulk insert.
	var insert bytes.Buffer
	if _, err := insert.WriteString(
		fmt.Sprintf(`INSERT INTO t.%s VALUES (%d, %d)`, name, 0, numRows-1)); err != nil {
		return err
	}
	for i := 1; i < numRows; i++ {
		if _, err := insert.WriteString(fmt.Sprintf(` ,(%d, %d)`, i, numRows-i)); err != nil {
			return err
		}
	}
	_, err := sqlDB.Exec(insert.String())
	return err
}
