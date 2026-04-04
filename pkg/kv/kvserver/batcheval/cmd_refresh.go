// Copyright 2017 The Cockroach Authors.
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

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
)

func init() {
	RegisterReadOnlyCommand(roachpb.Refresh, DefaultDeclareKeys, Refresh)
}

// Refresh checks whether the key has any values written in the interval
// (args.RefreshFrom, header.Timestamp].
func Refresh(
	ctx context.Context, reader storage.Reader, cArgs CommandArgs, resp roachpb.Response,
) (result.Result, error) {
	args := cArgs.Args.(*roachpb.RefreshRequest)
	h := cArgs.Header

	if h.Txn == nil {
		return result.Result{}, errors.AssertionFailedf("no transaction specified to %s", args.Method())
	}

	// We're going to refresh up to the transaction's read timestamp.
	if h.Timestamp != h.Txn.WriteTimestamp {
		// We're expecting the read and write timestamp to have converged before the
		// Refresh request was sent.
		log.Fatalf(ctx, "expected provisional commit ts %s == read ts %s. txn: %s", h.Timestamp,
			h.Txn.WriteTimestamp, h.Txn)
	}
	refreshTo := h.Timestamp

	refreshFrom := args.RefreshFrom
	if refreshFrom.IsEmpty() {
		return result.Result{}, errors.AssertionFailedf("empty RefreshFrom: %s", args)
	}

	// Get the most recent committed value and return any intent by
	// specifying consistent=false. Note that we include tombstones,
	// which must be considered as updates on refresh.
	log.VEventf(ctx, 2, "refresh %s @[%s-%s]", args.Span(), refreshFrom, refreshTo)
	val, intent, err := storage.MVCCGet(ctx, reader, args.Key, refreshTo, storage.MVCCGetOptions{
		Inconsistent: true,
		Tombstones:   true,
	})

	if err != nil {
		return result.Result{}, err
	} else if val != nil {
		if ts := val.Timestamp; refreshFrom.Less(ts) {
			return result.Result{},
				roachpb.NewRefreshFailedError(roachpb.RefreshFailedError_REASON_COMMITTED_VALUE, args.Key, ts)
		}
	}

	// Check if an intent which is not owned by this transaction was written
	// at or beneath the refresh timestamp.
	if intent != nil && intent.Txn.ID != h.Txn.ID {
		return result.Result{}, roachpb.NewRefreshFailedError(roachpb.RefreshFailedError_REASON_INTENT,
			intent.Key, intent.Txn.WriteTimestamp)
	}

	return result.Result{}, nil
}
