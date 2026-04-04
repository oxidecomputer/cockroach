// Copyright 2021 The Cockroach Authors.
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

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/spanset"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
)

func init() {
	RegisterReadWriteCommand(roachpb.Barrier, declareKeysBarrier, Barrier)
}

func declareKeysBarrier(
	_ ImmutableRangeState,
	_ *roachpb.Header,
	req roachpb.Request,
	latchSpans, _ *spanset.SpanSet,
	_ time.Duration,
) {
	// Barrier is special-cased in the concurrency manager to *not* actually
	// grab these latches. Instead, any conflicting latches with these are waited
	// on, but new latches aren't inserted.
	//
	// This will conflict with all writes and reads over the span, regardless of
	// timestamp. Note that this guarantees that concurrent writes will be flushed
	// out before this request is evaluated, but there is no guarantee regarding
	// flushing out of concurrent reads since they could be getting evaluated on a
	// follower. We don't currently need any guarantees regarding concurrent
	// reads, so this is acceptable.
	latchSpans.AddNonMVCC(spanset.SpanReadWrite, req.Header().Span())
}

// Barrier evaluation is a no-op, as all the latch waiting happens in
// the latch manager.
func Barrier(
	_ context.Context, _ storage.ReadWriter, cArgs CommandArgs, response roachpb.Response,
) (result.Result, error) {
	resp := response.(*roachpb.BarrierResponse)
	resp.Timestamp = cArgs.EvalCtx.Clock().Now()

	return result.Result{}, nil
}
