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

package acceptance

import (
	"context"
	gosql "database/sql"
	"net/http"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/acceptance/cluster"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestDebugRemote(t *testing.T) {
	s := log.Scope(t)
	defer s.Close(t)
	// TODO(knz): This test can probably move to a unit test inside the
	// server package
	RunDocker(t, testDebugRemote)
}

func testDebugRemote(t *testing.T) {
	cfg := cluster.TestConfig{
		Name:     "TestDebugRemote",
		Duration: *flagDuration,
		Nodes:    []cluster.NodeConfig{{Stores: []cluster.StoreConfig{{}}}},
	}
	ctx := context.Background()
	l := StartCluster(ctx, t, cfg).(*cluster.DockerCluster)
	defer l.AssertAndStop(ctx, t)

	db, err := gosql.Open("postgres", l.PGUrl(ctx, 0))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	stdout, stderr, err := l.ExecCLI(ctx, 0, []string{"auth-session", "login", "root", "--only-cookie"})
	if err != nil {
		t.Fatalf("auth-session failed: %s\nstdout: %s\nstderr: %s\n", err, stdout, stderr)
	}
	cookie := strings.Trim(stdout, "\n")

	for i, url := range []string{
		"/debug/",
		"/debug/pprof",
		"/debug/requests",
		"/debug/range?id=1",
		"/debug/certificates",
		"/debug/logspy?duration=1ns",
	} {
		t.Run(url, func(t *testing.T) {
			req, err := http.NewRequest("GET", l.URL(ctx, 0)+url, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Cookie", cookie)
			resp, err := cluster.HTTPClient.Do(req)
			if err != nil {
				t.Fatalf("%d: %v", i, err)
			}
			resp.Body.Close()

			if http.StatusOK != resp.StatusCode {
				t.Fatalf("%d: expected %d, but got %d", i, http.StatusOK, resp.StatusCode)
			}
		})
	}
}
