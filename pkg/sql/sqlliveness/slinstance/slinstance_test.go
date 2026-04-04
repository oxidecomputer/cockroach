// Copyright 2020 The Cockroach Authors.
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

package slinstance_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness/slinstance"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlliveness/slstorage"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
	"github.com/stretchr/testify/require"
)

func TestSQLInstance(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx, stopper := context.Background(), stop.NewStopper()
	defer stopper.Stop(ctx)

	mClock := hlc.NewManualClock(hlc.UnixNano())
	clock := hlc.NewClock(mClock.UnixNano, time.Nanosecond)
	settings := cluster.MakeTestingClusterSettingsWithVersions(
		clusterversion.TestingBinaryVersion,
		clusterversion.TestingBinaryMinSupportedVersion,
		true /* initializeVersion */)
	slinstance.DefaultTTL.Override(ctx, &settings.SV, 2*time.Microsecond)
	slinstance.DefaultHeartBeat.Override(ctx, &settings.SV, time.Microsecond)

	fakeStorage := slstorage.NewFakeStorage()
	sqlInstance := slinstance.NewSQLInstance(stopper, clock, fakeStorage, settings, nil)
	sqlInstance.Start(ctx)

	// Add one more instance to introduce concurrent access to storage.
	dummy := slinstance.NewSQLInstance(stopper, clock, fakeStorage, settings, nil)
	dummy.Start(ctx)

	s1, err := sqlInstance.Session(ctx)
	require.NoError(t, err)
	a, err := fakeStorage.IsAlive(ctx, s1.ID())
	require.NoError(t, err)
	require.True(t, a)

	s2, err := sqlInstance.Session(ctx)
	require.NoError(t, err)
	require.Equal(t, s1.ID(), s2.ID())

	_ = fakeStorage.Delete(ctx, s2.ID())
	t.Logf("deleted session %s", s2.ID())
	a, err = fakeStorage.IsAlive(ctx, s2.ID())
	require.NoError(t, err)
	require.False(t, a)

	var s3 sqlliveness.Session
	require.Eventually(
		t,
		func() bool {
			s3, err = sqlInstance.Session(ctx)
			if err != nil {
				t.Fatal(err)
			}
			return s3.ID().String() != s2.ID().String()
		},
		time.Second, 10*time.Millisecond,
	)

	a, err = fakeStorage.IsAlive(ctx, s3.ID())
	require.NoError(t, err)
	require.True(t, a)
	require.NotEqual(t, s2.ID(), s3.ID())

	// Force next call to Session to fail.
	stopper.Stop(ctx)
	sqlInstance.ClearSessionForTest(ctx)
	_, err = sqlInstance.Session(ctx)
	require.Error(t, err)
}
