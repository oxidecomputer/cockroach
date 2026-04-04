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

package readsummary

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/readsummary/rspb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/storage/enginepb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
)

// Load loads the range's prior read summary. The function returns a nil summary
// if one does not already exist.
func Load(
	ctx context.Context, reader storage.Reader, rangeID roachpb.RangeID,
) (*rspb.ReadSummary, error) {
	var sum rspb.ReadSummary
	key := keys.RangePriorReadSummaryKey(rangeID)
	found, err := storage.MVCCGetProto(ctx, reader, key, hlc.Timestamp{}, &sum, storage.MVCCGetOptions{})
	if !found {
		return nil, err
	}
	return &sum, err
}

// Set persists a range's prior read summary.
func Set(
	ctx context.Context,
	readWriter storage.ReadWriter,
	rangeID roachpb.RangeID,
	ms *enginepb.MVCCStats,
	sum *rspb.ReadSummary,
) error {
	key := keys.RangePriorReadSummaryKey(rangeID)
	return storage.MVCCPutProto(ctx, readWriter, ms, key, hlc.Timestamp{}, nil, sum)
}
