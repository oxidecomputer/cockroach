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
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/errors"
)

func init() {
	RegisterReadOnlyCommand(roachpb.QueryTxn, declareKeysQueryTransaction, QueryTxn)
}

func declareKeysQueryTransaction(
	_ ImmutableRangeState,
	_ *roachpb.Header,
	req roachpb.Request,
	latchSpans, _ *spanset.SpanSet,
	_ time.Duration,
) {
	qr := req.(*roachpb.QueryTxnRequest)
	latchSpans.AddNonMVCC(spanset.SpanReadOnly, roachpb.Span{Key: keys.TransactionKey(qr.Txn.Key, qr.Txn.ID)})
}

// QueryTxn fetches the current state of a transaction.
// This method is used to continually update the state of a txn
// which is blocked waiting to resolve a conflicting intent. It
// fetches the complete transaction record to determine whether
// priority or status has changed and also fetches a list of
// other txns which are waiting on this transaction in order
// to find dependency cycles.
func QueryTxn(
	ctx context.Context, reader storage.Reader, cArgs CommandArgs, resp roachpb.Response,
) (result.Result, error) {
	args := cArgs.Args.(*roachpb.QueryTxnRequest)
	h := cArgs.Header
	reply := resp.(*roachpb.QueryTxnResponse)

	if h.Txn != nil {
		return result.Result{}, ErrTransactionUnsupported
	}
	if h.WriteTimestamp().Less(args.Txn.MinTimestamp) {
		// This condition must hold for the timestamp cache access in
		// SynthesizeTxnFromMeta to be safe.
		return result.Result{}, errors.AssertionFailedf("QueryTxn request timestamp %s less than txn MinTimestamp %s",
			h.Timestamp, args.Txn.MinTimestamp)
	}
	if !args.Key.Equal(args.Txn.Key) {
		return result.Result{}, errors.AssertionFailedf("QueryTxn request key %s does not match txn key %s",
			args.Key, args.Txn.Key)
	}
	key := keys.TransactionKey(args.Txn.Key, args.Txn.ID)

	// Fetch transaction record; if missing, attempt to synthesize one.
	ok, err := storage.MVCCGetProto(
		ctx, reader, key, hlc.Timestamp{}, &reply.QueriedTxn, storage.MVCCGetOptions{},
	)
	if err != nil {
		return result.Result{}, err
	}
	if ok {
		reply.TxnRecordExists = true
	} else {
		// The transaction hasn't written a transaction record yet.
		// Attempt to synthesize it from the provided TxnMeta.
		reply.QueriedTxn = SynthesizeTxnFromMeta(ctx, cArgs.EvalCtx, args.Txn)
	}

	// Get the list of txns waiting on this txn.
	reply.WaitingTxns = cArgs.EvalCtx.GetConcurrencyManager().GetDependents(args.Txn.ID)
	return result.Result{}, nil
}
