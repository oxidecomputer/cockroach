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

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/spanset"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
)

func init() {
	RegisterReadWriteCommand(roachpb.ResolveIntentRange, declareKeysResolveIntentRange, ResolveIntentRange)
}

func declareKeysResolveIntentRange(
	rs ImmutableRangeState,
	_ *roachpb.Header,
	req roachpb.Request,
	latchSpans, _ *spanset.SpanSet,
	_ time.Duration,
) {
	declareKeysResolveIntentCombined(rs, req, latchSpans)
}

// ResolveIntentRange resolves write intents in the specified
// key range according to the status of the transaction which created it.
func ResolveIntentRange(
	ctx context.Context, readWriter storage.ReadWriter, cArgs CommandArgs, resp roachpb.Response,
) (result.Result, error) {
	args := cArgs.Args.(*roachpb.ResolveIntentRangeRequest)
	h := cArgs.Header
	ms := cArgs.Stats

	if h.Txn != nil {
		return result.Result{}, ErrTransactionUnsupported
	}

	update := args.AsLockUpdate()
	numKeys, resumeSpan, err := storage.MVCCResolveWriteIntentRange(
		ctx, readWriter, ms, update, h.MaxSpanRequestKeys)
	if err != nil {
		return result.Result{}, err
	}
	reply := resp.(*roachpb.ResolveIntentRangeResponse)
	reply.NumKeys = numKeys
	if resumeSpan != nil {
		update.EndKey = resumeSpan.Key
		reply.ResumeSpan = resumeSpan
		// The given MaxSpanRequestKeys really specifies the number of intents
		// resolved, not the number of keys scanned. We could return
		// RESUME_INTENT_LIMIT here, but since the given limit is a key limit we
		// return RESUME_KEY_LIMIT for symmetry.
		reply.ResumeReason = roachpb.RESUME_KEY_LIMIT
	}

	var res result.Result
	res.Local.ResolvedLocks = []roachpb.LockUpdate{update}
	res.Local.Metrics = resolveToMetricType(args.Status, args.Poison)

	if WriteAbortSpanOnResolve(args.Status, args.Poison, numKeys > 0) {
		if err := UpdateAbortSpan(ctx, cArgs.EvalCtx, readWriter, ms, args.IntentTxn, args.Poison); err != nil {
			return result.Result{}, err
		}
	}
	return res, nil
}
