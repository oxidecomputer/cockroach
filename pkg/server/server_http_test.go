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
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestHSTS(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	s, db, _ := serverutils.StartServer(t, base.TestServerArgs{})
	defer s.Stopper().Stop(ctx)

	httpClient, err := s.GetHTTPClient()
	require.NoError(t, err)
	defer httpClient.CloseIdleConnections()

	secureClient, err := s.GetAuthenticatedHTTPClient(false)
	require.NoError(t, err)
	defer secureClient.CloseIdleConnections()

	urlsToTest := []string{"/", "/_status/vars", "/index.html"}

	adminURLHTTPS := s.AdminURL()
	adminURLHTTP := strings.Replace(adminURLHTTPS, "https", "http", 1)

	for _, u := range urlsToTest {
		resp, err := httpClient.Get(adminURLHTTP + u)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Empty(t, resp.Header.Get(hstsHeaderKey))

		resp, err = secureClient.Get(adminURLHTTPS + u)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Empty(t, resp.Header.Get(hstsHeaderKey))
	}

	_, err = db.Exec("SET cluster setting server.hsts.enabled = true")
	require.NoError(t, err)

	for _, u := range urlsToTest {
		resp, err := httpClient.Get(adminURLHTTP + u)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, resp.Header.Get(hstsHeaderKey), hstsHeaderValue)

		resp, err = secureClient.Get(adminURLHTTPS + u)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, resp.Header.Get(hstsHeaderKey), hstsHeaderValue)
	}
	_, err = db.Exec("SET cluster setting server.hsts.enabled = false")
	require.NoError(t, err)

	for _, u := range urlsToTest {
		resp, err := httpClient.Get(adminURLHTTP + u)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Empty(t, resp.Header.Get(hstsHeaderKey))

		resp, err = secureClient.Get(adminURLHTTPS + u)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Empty(t, resp.Header.Get(hstsHeaderKey))
	}
}
