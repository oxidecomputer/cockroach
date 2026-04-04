// Copyright 2023 The Cockroach Authors.
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

//go:build linux
// +build linux

package cli

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestFindOpenSocketFDs(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(500 * time.Millisecond)
		rw.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = s.Client().Get(s.URL)
	}()
	defer wg.Wait()

	fds, _ := findOpenSocketFDs()
	require.Greater(t, len(fds), 0)
	s.CloseClientConnections()
}
