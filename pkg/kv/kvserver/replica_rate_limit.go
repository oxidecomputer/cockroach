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

	"github.com/cockroachdb/cockroach/pkg/multitenant/tenantcostmodel"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// maybeRateLimitBatch may block the batch waiting to be rate-limited. Note that
// the replica must be initialized and thus there is no synchronization issue
// on the tenantRateLimiter.
func (r *Replica) maybeRateLimitBatch(ctx context.Context, ba *roachpb.BatchRequest) error {
	if r.tenantLimiter == nil {
		return nil
	}
	tenantID, ok := roachpb.TenantFromContext(ctx)
	if !ok || tenantID == roachpb.SystemTenantID {
		return nil
	}
	return r.tenantLimiter.Wait(ctx, tenantcostmodel.MakeRequestInfo(ba, 1))
}

// recordImpactOnRateLimiter is used to record a read against the tenant rate
// limiter.
func (r *Replica) recordImpactOnRateLimiter(
	ctx context.Context, br *roachpb.BatchResponse, isReadOnly bool,
) {
	if r.tenantLimiter == nil || br == nil || !isReadOnly {
		return
	}

	r.tenantLimiter.RecordRead(ctx, tenantcostmodel.MakeResponseInfo(br, isReadOnly))
}
