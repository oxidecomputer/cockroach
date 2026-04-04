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

package clisqlclient

import (
	"context"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

// sqlTxnShim implements the crdb.Tx interface.
//
// It exists to support crdb.ExecuteInTxn. Normally, we'd hand crdb.ExecuteInTxn
// a sql.Txn, but sqlConn predates go1.8's support for multiple result sets and
// so deals directly with the lib/pq driver. See #14964.
//
// TODO(knz): This code is incorrect, see
// https://github.com/cockroachdb/cockroach/issues/67261
type sqlTxnShim struct {
	conn *sqlConn
}

var _ crdb.Tx = sqlTxnShim{}

func (t sqlTxnShim) Commit(ctx context.Context) error {
	return t.conn.Exec(ctx, `COMMIT`)
}

func (t sqlTxnShim) Rollback(ctx context.Context) error {
	return t.conn.Exec(ctx, `ROLLBACK`)
}

func (t sqlTxnShim) Exec(ctx context.Context, query string, values ...interface{}) error {
	if len(values) != 0 {
		panic("sqlTxnShim.ExecContext must not be called with values")
	}
	return t.conn.Exec(ctx, query)
}
