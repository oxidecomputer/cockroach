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

package clierror_test

import (
	"context"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/cli"
	"github.com/cockroachdb/cockroach/pkg/cli/clierror"
	"github.com/cockroachdb/cockroach/pkg/cli/clisqlclient"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
)

// This test checks that IsSQLSyntaxError works. It could stop working if e.g.
// the surrounding code stops using lib/pq as SQL driver, and/or the error type
// from query execution is not pq.Error any more.
func TestIsSQLSyntaxError(t *testing.T) {
	defer leaktest.AfterTest(t)()

	p := cli.TestCLIParams{T: t}
	c := cli.NewCLITest(p)
	defer c.Cleanup()

	url, cleanup := sqlutils.PGUrl(t, c.ServingSQLAddr(), t.Name(), url.User(security.RootUser))
	defer cleanup()

	var sqlConnCtx clisqlclient.Context
	conn := sqlConnCtx.MakeSQLConn(ioutil.Discard, ioutil.Discard, url.String())
	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	_, err := conn.QueryRow(context.Background(), `INVALID SYNTAX`)
	if !clierror.IsSQLSyntaxError(err) {
		t.Fatalf("expected error to be recognized as syntax error: %+v", err)
	}
}
