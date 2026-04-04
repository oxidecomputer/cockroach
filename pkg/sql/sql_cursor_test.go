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

package sql

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// Make sure that preparing cursor statements don't cause problems.
func TestPrepareCursors(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	srv, db, _ := serverutils.StartServer(t, base.TestServerArgs{Insecure: true})
	defer srv.Stopper().Stop(context.Background())
	defer db.Close()

	ctx := context.Background()
	conn, err := db.Conn(ctx)
	require.NoError(t, err)

	t.Run("prepare_declare_raw_txn", func(t *testing.T) {
		// Make sure that preparing a DECLARE defers errors until execution.
		stmt, err := conn.PrepareContext(ctx, "DECLARE foo CURSOR FOR VALUES (1), (2)")
		require.NoError(t, err)

		_, err = stmt.Exec()
		require.EqualError(t, err, "pq: DECLARE CURSOR can only be used in transaction blocks")

		// Make sure that we can use our prepared statement from before to
		// successfully execute a declare cursor within a transaction.
		// We need to execute a raw BEGIN so that we can reuse our pre-prepared txn.
		_, err = conn.ExecContext(ctx, "BEGIN TRANSACTION")
		require.NoError(t, err)

		_, err = stmt.Exec()
		require.NoError(t, err)

		stmt, err = conn.PrepareContext(ctx, "FETCH 2 foo")
		require.NoError(t, err)
		r, err := stmt.Query()
		require.NoError(t, err)
		var actual int
		r.Next()
		require.NoError(t, r.Scan(&actual))
		require.Equal(t, 1, actual)
		more := r.Next()
		require.Equal(t, true, more)
		require.NoError(t, r.Scan(&actual))
		require.Equal(t, 2, actual)
		more = r.Next()
		require.Equal(t, false, more)

		stmt, err = conn.PrepareContext(ctx, "MOVE 1 foo")
		require.NoError(t, err)
		_, err = stmt.Exec()
		require.NoError(t, err)

		_, err = conn.ExecContext(ctx, "COMMIT")
		require.NoError(t, err)
	})

	t.Run("prepare_declare_driver_txn", func(t *testing.T) {
		// Make sure that we can use the driver-level txn support to do the same thing.
		tx, err := conn.BeginTx(context.Background(), nil /* opts */)
		require.NoError(t, err)
		stmt, err := tx.Prepare("DECLARE foo CURSOR FOR VALUES (1), (2)")
		require.NoError(t, err)
		_, err = stmt.Exec()
		require.NoError(t, err)

		stmt, err = tx.Prepare("FETCH 2 foo")
		require.NoError(t, err)
		r, err := stmt.Query()
		require.NoError(t, err)
		var actual int
		r.Next()
		require.NoError(t, r.Scan(&actual))
		require.Equal(t, 1, actual)
		more := r.Next()
		require.Equal(t, true, more)
		require.NoError(t, r.Scan(&actual))
		require.Equal(t, 2, actual)
		more = r.Next()
		require.Equal(t, false, more)

		stmt, err = tx.Prepare("MOVE 1 foo")
		require.NoError(t, err)
		_, err = stmt.Exec()
		require.NoError(t, err)

		require.NoError(t, tx.Commit())
	})

	// Make sure that we can use the automatic prepare support (when sending
	// placeholders) to do the same thing.
	t.Run("prepare_declare_placeholder", func(t *testing.T) {
		_, err = conn.ExecContext(ctx, "BEGIN TRANSACTION")
		require.NoError(t, err)

		_, err = conn.ExecContext(ctx, "DECLARE foo CURSOR FOR SELECT 1 WHERE $1", true)
		require.NoError(t, err)

		stmt, err := conn.PrepareContext(ctx, "FETCH 1 foo")
		require.NoError(t, err)
		r, err := stmt.Query()
		require.NoError(t, err)
		var actual int
		r.Next()
		require.NoError(t, r.Scan(&actual))
		require.Equal(t, 1, actual)
		more := r.Next()
		require.Equal(t, false, more)
	})
}
