// Copyright 2016 The Cockroach Authors.
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

package kvserver

import (
	"fmt"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestReplicaGCShouldQueue(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ts := func(t time.Duration) hlc.Timestamp {
		return hlc.Timestamp{WallTime: t.Nanoseconds()}
	}

	base := 2 * (ReplicaGCQueueSuspectCheckInterval + ReplicaGCQueueCheckInterval)
	var (
		z       = ts(0)
		bTS     = ts(base)
		sTS     = ts(base + ReplicaGCQueueSuspectCheckInterval)
		sTSnext = ts(base + ReplicaGCQueueSuspectCheckInterval + 1)
		iTS     = ts(base + ReplicaGCQueueCheckInterval)
		iTSnext = ts(base + ReplicaGCQueueCheckInterval + 1)
	)

	testcases := []struct {
		now            hlc.Timestamp
		lastCheck      hlc.Timestamp
		isSuspect      bool
		expectQueue    bool
		expectPriority float64
	}{
		// All timestamps current: suspect plays no role.
		{now: z, lastCheck: z, isSuspect: true, expectQueue: false, expectPriority: 0},
		// Threshold: no action taken.
		{now: sTS, lastCheck: bTS, isSuspect: true, expectQueue: false, expectPriority: 0},
		// Last processed recently: suspect still gets processed eagerly.
		{now: sTSnext, lastCheck: bTS, isSuspect: true, expectQueue: true, expectPriority: 1},
		// Last processed recently: non-suspect stays put.
		{now: sTSnext, lastCheck: bTS, isSuspect: false, expectQueue: false, expectPriority: 0},
		// No effect until iTS crossed.
		{now: iTS, lastCheck: bTS, isSuspect: false, expectQueue: false, expectPriority: 0},
		{now: iTSnext, lastCheck: bTS, isSuspect: false, expectQueue: true, expectPriority: 0},
		// Verify again that suspect increases priority.
		{now: iTSnext, lastCheck: bTS, isSuspect: true, expectQueue: true, expectPriority: 1},
	}
	for _, tc := range testcases {
		tc := tc
		name := fmt.Sprintf("now=%s lastCheck=%s isSuspect=%v", tc.now, tc.lastCheck, tc.isSuspect)
		t.Run(name, func(t *testing.T) {
			shouldQueue, priority := replicaGCShouldQueueImpl(tc.now, tc.lastCheck, tc.isSuspect)
			require.Equal(t, tc.expectQueue, shouldQueue)
			require.Equal(t, tc.expectPriority, priority)
		})
	}
}
