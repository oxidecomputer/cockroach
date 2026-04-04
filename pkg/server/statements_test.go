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

package server

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/rpc"
	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/sql/tests"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// TestStatements ensures that the Statements endpoint is accessible
// via gRPC and returns information reflecting recently run SQL
// queries.
func TestStatements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()

	params, _ := tests.CreateTestServerParams()
	testServer, db, _ := serverutils.StartServer(t, params)
	defer testServer.Stopper().Stop(ctx)

	conn, err := testServer.RPCContext().GRPCDialNode(
		testServer.RPCAddr(), testServer.NodeID(), rpc.DefaultClass,
	).Connect(ctx)
	require.NoError(t, err)

	client := serverpb.NewStatusClient(conn)

	testQuery := "CREATE TABLE foo (id INT8)"
	_, err = db.Exec(testQuery)
	require.NoError(t, err)

	resp, err := client.Statements(ctx, &serverpb.StatementsRequest{NodeID: "local"})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Statements)
	require.NotEmpty(t, resp.Transactions)

	queries := make([]string, len(resp.Statements))
	for _, s := range resp.Statements {
		queries = append(queries, s.Key.KeyData.Query)
	}
	require.Contains(t, queries, testQuery)
}

func TestStatementsExcludeStats(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()

	params, _ := tests.CreateTestServerParams()
	testServer, db, _ := serverutils.StartServer(t, params)
	defer testServer.Stopper().Stop(ctx)

	conn, err := testServer.RPCContext().GRPCDialNode(
		testServer.RPCAddr(), testServer.NodeID(), rpc.DefaultClass,
	).Connect(ctx)
	require.NoError(t, err)

	client := serverpb.NewStatusClient(conn)

	testQuery := "CREATE TABLE foo (id INT8)"
	_, err = db.Exec(testQuery)
	require.NoError(t, err)

	t.Run("exclude-statements", func(t *testing.T) {
		resp, err := client.Statements(ctx, &serverpb.StatementsRequest{
			NodeID:    "local",
			FetchMode: serverpb.StatementsRequest_TxnStatsOnly,
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.Transactions)
		require.Empty(t, resp.Statements)
	})

	t.Run("exclude-transactions", func(t *testing.T) {
		resp, err := client.Statements(ctx, &serverpb.StatementsRequest{
			NodeID:    "local",
			FetchMode: serverpb.StatementsRequest_StmtStatsOnly,
		})
		require.NoError(t, err)
		require.Empty(t, resp.Transactions)
		require.NotEmpty(t, resp.Statements)
	})
}
