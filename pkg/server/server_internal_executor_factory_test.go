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
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestInternalExecutorClearsMonitorMemory(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	s, _, _ := serverutils.StartServer(t, base.TestServerArgs{})
	defer s.Stopper().Stop(ctx)

	mon := s.(*TestServer).sqlServer.internalExecutorFactoryMemMonitor
	ief := s.ExecutorConfig().(sql.ExecutorConfig).InternalExecutorFactory
	sessionData := sql.NewFakeSessionData(&s.ClusterSettings().SV)
	ie := ief(ctx, sessionData)
	rows, err := ie.QueryIteratorEx(ctx, "test", nil, sessiondata.NodeUserSessionDataOverride, `SELECT 1`)
	require.NoError(t, err)
	require.Greater(t, mon.AllocBytes(), int64(0))
	err = rows.Close()
	require.NoError(t, err)
	s.Stopper().Stop(ctx)
	require.Equal(t, mon.AllocBytes(), int64(0))
}
