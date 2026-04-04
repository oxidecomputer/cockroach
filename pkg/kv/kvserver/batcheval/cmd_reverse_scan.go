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
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/concurrency/lock"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
)

func init() {
	RegisterReadOnlyCommand(roachpb.ReverseScan, DefaultDeclareIsolatedKeys, ReverseScan)
}

// ReverseScan scans the key range specified by start key through
// end key in descending order up to some maximum number of results.
// maxKeys stores the number of scan results remaining for this batch
// (MaxInt64 for no limit).
func ReverseScan(
	ctx context.Context, reader storage.Reader, cArgs CommandArgs, resp roachpb.Response,
) (result.Result, error) {
	args := cArgs.Args.(*roachpb.ReverseScanRequest)
	h := cArgs.Header
	reply := resp.(*roachpb.ReverseScanResponse)

	var res result.Result
	var scanRes storage.MVCCScanResult
	var err error

	avoidExcess := cArgs.EvalCtx.ClusterSettings().Version.IsActive(ctx,
		clusterversion.TargetBytesAvoidExcess)
	opts := storage.MVCCScanOptions{
		Inconsistent:           h.ReadConsistency != roachpb.CONSISTENT,
		Txn:                    h.Txn,
		Uncertainty:            cArgs.Uncertainty,
		MaxKeys:                h.MaxSpanRequestKeys,
		MaxIntents:             storage.MaxIntentsPerWriteIntentError.Get(&cArgs.EvalCtx.ClusterSettings().SV),
		TargetBytes:            h.TargetBytes,
		TargetBytesAvoidExcess: h.AllowEmpty || avoidExcess, // AllowEmpty takes precedence
		AllowEmpty:             h.AllowEmpty,
		WholeRowsOfSize:        h.WholeRowsOfSize,
		FailOnMoreRecent:       args.KeyLocking != lock.None,
		Reverse:                true,
		MemoryAccount:          cArgs.EvalCtx.GetResponseMemoryAccount(),
	}

	switch args.ScanFormat {
	case roachpb.BATCH_RESPONSE:
		scanRes, err = storage.MVCCScanToBytes(
			ctx, reader, args.Key, args.EndKey, h.Timestamp, opts)
		if err != nil {
			return result.Result{}, err
		}
		reply.BatchResponses = scanRes.KVData
	case roachpb.KEY_VALUES:
		scanRes, err = storage.MVCCScan(
			ctx, reader, args.Key, args.EndKey, h.Timestamp, opts)
		if err != nil {
			return result.Result{}, err
		}
		reply.Rows = scanRes.KVs
	default:
		panic(fmt.Sprintf("Unknown scanFormat %d", args.ScanFormat))
	}

	reply.NumKeys = scanRes.NumKeys
	reply.NumBytes = scanRes.NumBytes

	if scanRes.ResumeSpan != nil {
		reply.ResumeSpan = scanRes.ResumeSpan
		reply.ResumeReason = scanRes.ResumeReason
		reply.ResumeNextBytes = scanRes.ResumeNextBytes
	}

	if h.ReadConsistency == roachpb.READ_UNCOMMITTED {
		// NOTE: MVCCScan doesn't use a Prefix iterator, so we don't want to use
		// one in CollectIntentRows either so that we're guaranteed to use the
		// same cached iterator and observe a consistent snapshot of the engine.
		const usePrefixIter = false
		reply.IntentRows, err = CollectIntentRows(ctx, reader, usePrefixIter, scanRes.Intents)
		if err != nil {
			return result.Result{}, err
		}
	}

	if args.KeyLocking != lock.None && h.Txn != nil {
		err = acquireUnreplicatedLocksOnKeys(&res, h.Txn, args.ScanFormat, &scanRes)
		if err != nil {
			return result.Result{}, err
		}
	}
	res.Local.EncounteredIntents = scanRes.Intents
	return res, nil
}
