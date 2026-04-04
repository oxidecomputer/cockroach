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

package kvprober

import (
	"context"
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/stretchr/testify/require"
)

func TestQuarantinePool(t *testing.T) {
	ctx := context.Background()

	t.Run("disabled by default", func(t *testing.T) {
		m := &mock{
			t:      t,
			noPlan: true,
		}
		p := initTestProber(ctx, m)
		p.quarantineProbe(ctx, m)

		require.Zero(t, p.Metrics().WriteProbeQuarantineOldestDuration.Value())
		require.Zero(t, p.Metrics().WriteProbeAttempts.Count())
	})

	t.Run("add then remove from quarantine pool", func(t *testing.T) {
		// adds to the quarantine pool
		m := &mock{t: t, write: true, qWrite: true, writeErr: fmt.Errorf("inject write failure")}
		p := initTestProber(ctx, m)
		p.writeProbeImpl(ctx, m, m, m)

		require.Equal(t, 1, len(p.quarantineWritePool.mu.steps))
		require.Equal(t, 1, len(p.quarantineWritePool.mu.entryTimeMap))
		require.Equal(t, int64(1), p.Metrics().WriteProbeFailures.Count())

		// removes from the quarantine pool
		m = &mock{t: t, write: true, qWrite: true}
		p.writeProbeImpl(ctx, m, m, m)

		require.Zero(t, len(p.quarantineWritePool.mu.steps))
		require.Empty(t, p.quarantineWritePool.mu.entryTimeMap)
		require.Equal(t, int64(1), p.Metrics().WriteProbeFailures.Count())
	})

	t.Run("test maybeAdd false cases", func(t *testing.T) {
		m := &mock{t: t, write: true, qWrite: true, writeErr: fmt.Errorf("inject write failure")}
		p := initTestProber(ctx, m)
		p.writeProbeImpl(ctx, m, m, m)

		step, err := p.quarantineWritePool.next(ctx)
		require.NoError(t, err)

		// already in pool case
		added := p.quarantineWritePool.maybeAdd(ctx, step)
		require.False(t, added)

		// pool full case
		quarantinePoolSize.Override(ctx, &p.settings.SV, 1)
		mockStep := Step{RangeID: 1, Key: keys.LocalMax}
		added = p.quarantineWritePool.maybeAdd(ctx, mockStep)
		require.False(t, added)
	})

	t.Run("ensure empty q pool does not cause planning failure", func(t *testing.T) {
		m := &mock{t: t, write: true, qWrite: true, emptyQPool: true}
		p := initTestProber(ctx, m)
		p.writeProbeImpl(ctx, m, m, m)

		require.Zero(t, p.Metrics().ProbePlanFailures.Count())
	})
}
