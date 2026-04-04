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

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/datadriven"
	"github.com/stretchr/testify/require"
)

func TestExecSQL(t *testing.T) {
	defer leaktest.AfterTest(t)()
	sqlAPIClock = timeutil.NewManualTime(timeutil.FromUnixMicros(0))
	defer func() {
		sqlAPIClock = timeutil.DefaultTimeSource{}
	}()

	server, _, _ := serverutils.StartServer(t, base.TestServerArgs{})
	ctx := context.Background()
	defer server.Stopper().Stop(ctx)

	adminClient, err := server.GetAdminAuthenticatedHTTPClient()
	require.NoError(t, err)

	nonAdminClient, err := server.GetAuthenticatedHTTPClient(false)
	require.NoError(t, err)

	datadriven.RunTest(t, "testdata/api_v2_sql",
		func(t *testing.T, d *datadriven.TestData) string {
			if d.Cmd != "sql" {
				t.Fatal("Only sql command is accepted in this test")
			}

			var client http.Client
			if d.HasArg("admin") {
				client = adminClient
			}
			if d.HasArg("non-admin") {
				client = nonAdminClient
			}

			resp, err := client.Post(
				server.AdminURL()+"/api/v2/sql/", "application/json",
				bytes.NewReader([]byte(d.Input)),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			r, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if d.HasArg("expect-error") {
				type jsonError struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}
				type errorResp struct {
					Error jsonError `json:"error"`
				}

				er := errorResp{}
				err := json.Unmarshal(r, &er)
				require.NoError(t, err)
				return fmt.Sprintf("%s|%s", er.Error.Code, er.Error.Message)
			}
			var u interface{}
			err = json.Unmarshal(r, &u)
			require.NoError(t, err)
			s, err := json.MarshalIndent(u, "", " ")
			require.NoError(t, err)
			return string(s)
		},
	)
}
