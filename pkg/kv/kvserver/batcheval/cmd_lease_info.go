// Copyright 2014 The Cockroach Authors.
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

package batcheval

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/spanset"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
)

func init() {
	RegisterReadOnlyCommand(roachpb.LeaseInfo, declareKeysLeaseInfo, LeaseInfo)
}

func declareKeysLeaseInfo(
	rs ImmutableRangeState,
	_ *roachpb.Header,
	_ roachpb.Request,
	latchSpans, _ *spanset.SpanSet,
	_ time.Duration,
) {
	latchSpans.AddNonMVCC(spanset.SpanReadOnly, roachpb.Span{Key: keys.RangeLeaseKey(rs.GetRangeID())})
}

// LeaseInfo returns information about the lease holder for the range.
func LeaseInfo(
	ctx context.Context, reader storage.Reader, cArgs CommandArgs, resp roachpb.Response,
) (result.Result, error) {
	reply := resp.(*roachpb.LeaseInfoResponse)
	lease, nextLease := cArgs.EvalCtx.GetLease()
	if nextLease != (roachpb.Lease{}) {
		// If there's a lease request in progress, speculatively return that future
		// lease.
		reply.Lease = nextLease
		reply.CurrentLease = &lease
	} else {
		reply.Lease = lease
	}
	reply.EvaluatedBy = cArgs.EvalCtx.StoreID()
	return result.Result{}, nil
}
