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

package kvserver

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage/enginepb"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// TestTenantsStorageMetricsConcurrency exercises the concurrency logic of the
// TenantsStorageMetrics and ensures that none of the assertions are hit.
// The test doesn't meaningfully exercise the logic which is tested elsewhere.
func TestTenantsStorageMetricsConcurrency(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	const (
		tenants  = 3
		N        = 200
		maxSleep = time.Microsecond
		rounds   = 10
	)
	randDuration := func() time.Duration {
		return time.Duration(rand.Intn(int(maxSleep)))
	}

	var tenantIDs []roachpb.TenantID
	for id := uint64(1); id <= tenants; id++ {
		tenantIDs = append(tenantIDs, roachpb.MakeTenantID(id))
	}
	ctx := context.Background()
	sm := newTenantsStorageMetrics()
	// Launch N goroutines and have them all acquire a random tenant, then sleep
	// a random tiny amount, increment the metrics, then release. We want to
	// ensure that the refCount is never in an illegal state.
	run := func() {
		for i := 0; i < rounds; i++ {
			tid := tenantIDs[rand.Intn(tenants)]

			time.Sleep(randDuration())
			ref := sm.acquireTenant(tid)

			time.Sleep(randDuration())
			sm.incMVCCGauges(ctx, ref, enginepb.MVCCStats{})

			time.Sleep(randDuration())
			sm.releaseTenant(ctx, ref)
		}
	}
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() { defer wg.Done(); run() }()
	}
	wg.Wait()
}
