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

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/batcheval/result"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
)

func init() {
	RegisterReadWriteCommand(roachpb.InitPut, DefaultDeclareIsolatedKeys, InitPut)
}

// InitPut sets the value for a specified key only if it doesn't exist. It
// returns a ConditionFailedError if the key exists with an existing value that
// is different from the value provided. If FailOnTombstone is set to true,
// tombstones count as mismatched values and will cause a ConditionFailedError.
func InitPut(
	ctx context.Context, readWriter storage.ReadWriter, cArgs CommandArgs, resp roachpb.Response,
) (result.Result, error) {
	args := cArgs.Args.(*roachpb.InitPutRequest)
	h := cArgs.Header

	var err error
	if args.Blind {
		err = storage.MVCCBlindInitPut(ctx, readWriter, cArgs.Stats, args.Key, h.Timestamp, args.Value, args.FailOnTombstones, h.Txn)
	} else {
		err = storage.MVCCInitPut(ctx, readWriter, cArgs.Stats, args.Key, h.Timestamp, args.Value, args.FailOnTombstones, h.Txn)
	}
	// NB: even if MVCC returns an error, it may still have written an intent
	// into the batch. This allows callers to consume errors like WriteTooOld
	// without re-evaluating the batch. This behavior isn't particularly
	// desirable, but while it remains, we need to assume that an intent could
	// have been written even when an error is returned. This is harmless if the
	// error is not consumed by the caller because the result will be discarded.
	return result.FromAcquiredLocks(h.Txn, args.Key), err
}
