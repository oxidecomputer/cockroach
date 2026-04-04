// Copyright 2022 The Cockroach Authors.
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

package clisqlclient

import (
	"context"
	"database/sql/driver"
	"io"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type copyFromer interface {
	CopyData(ctx context.Context, line string) (r driver.Result, err error)
	Exec(v []driver.Value) (r driver.Result, err error)
	Close() error
}

// CopyFromState represents an in progress COPY FROM.
type CopyFromState struct {
	driver.Tx
	copyFromer
}

// BeginCopyFrom starts a COPY FROM query.
func BeginCopyFrom(ctx context.Context, conn Conn, query string) (*CopyFromState, error) {
	txn, err := conn.(*sqlConn).conn.(driver.ConnBeginTx).BeginTx(ctx, driver.TxOptions{})
	if err != nil {
		return nil, err
	}
	stmt, err := txn.(driver.Conn).Prepare(query)
	if err != nil {
		return nil, errors.CombineErrors(err, txn.Rollback())
	}
	return &CopyFromState{Tx: txn, copyFromer: stmt.(copyFromer)}, nil
}

// copyFromRows is a mock Rows interface for COPY results.
type copyFromRows struct {
	r driver.Result
}

func (c copyFromRows) Close() error {
	return nil
}

func (c copyFromRows) Columns() []string {
	return nil
}

func (c copyFromRows) ColumnTypeScanType(index int) reflect.Type {
	return nil
}

func (c copyFromRows) ColumnTypeDatabaseTypeName(index int) string {
	return ""
}

func (c copyFromRows) ColumnTypeNames() []string {
	return nil
}

func (c copyFromRows) Result() driver.Result {
	return c.r
}

func (c copyFromRows) Tag() string {
	return "COPY"
}

func (c copyFromRows) Next(values []driver.Value) error {
	return io.EOF
}

func (c copyFromRows) NextResultSet() (bool, error) {
	return false, nil
}

// Cancel cancels a COPY FROM query from completing.
func (c *CopyFromState) Cancel() error {
	return errors.CombineErrors(c.copyFromer.Close(), c.Tx.Rollback())
}

// Commit completes a COPY FROM query by committing lines to the database.
func (c *CopyFromState) Commit(ctx context.Context, cleanupFunc func(), lines string) QueryFn {
	return func(ctx context.Context, conn Conn) (Rows, bool, error) {
		defer cleanupFunc()
		rows, isMulti, err := func() (Rows, bool, error) {
			// Do not send anything if it is just an empty string.
			if lines != "" {
				for _, l := range strings.Split(lines, "\n") {
					_, err := c.copyFromer.CopyData(ctx, l)
					if err != nil {
						return nil, false, err
					}
				}
			}
			r, err := c.copyFromer.Exec(nil)
			if err != nil {
				return nil, false, err
			}
			return copyFromRows{r: r}, false, c.Tx.Commit()
		}()
		if err != nil {
			return rows, isMulti, errors.CombineErrors(err, errors.CombineErrors(c.copyFromer.Close(), c.Tx.Rollback()))
		}
		return rows, isMulti, err
	}
}
