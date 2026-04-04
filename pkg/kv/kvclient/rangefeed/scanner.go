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

package rangefeed

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/retry"
	"github.com/cockroachdb/cockroach/pkg/util/span"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
)

// runInitialScan will attempt to perform an initial data scan.
// It will retry in the face of errors and will only return upon
// success, context cancellation, or an error handling function which indicates
// that an error is unrecoverable. The return value will be true if the context
// was canceled or if the OnInitialScanError function indicated that the
// RangeFeed should stop.
func (f *RangeFeed) runInitialScan(
	ctx context.Context, n *log.EveryN, r *retry.Retry,
) (canceled bool) {
	onValue := func(kv roachpb.KeyValue) {
		v := roachpb.RangeFeedValue{
			Key:   kv.Key,
			Value: kv.Value,
		}

		// Mark the data as occurring at the initial timestamp, which is the
		// timestamp at which it was read.
		if !f.useRowTimestampInInitialScan {
			v.Value.Timestamp = f.initialTimestamp
		}

		// Supply the value from the scan as also the previous value to avoid
		// indicating that the value was previously deleted.
		if f.withDiff {
			v.PrevValue = v.Value
			v.PrevValue.Timestamp = hlc.Timestamp{}
		}

		// It's something of a bummer that we must allocate a new value for each
		// of these but the contract doesn't indicate that the value cannot be
		// retained so we have to assume that the callback may retain the value.
		f.onValue(ctx, &v)
	}

	getSpansToScan := f.getSpansToScan(ctx)

	r.Reset()
	for r.Next() {
		if err := f.client.Scan(ctx, getSpansToScan(), f.initialTimestamp, onValue, f.scanConfig); err != nil {
			if f.onInitialScanError != nil {
				if shouldStop := f.onInitialScanError(ctx, err); shouldStop {
					log.VEventf(ctx, 1, "stopping due to error: %v", err)
					return true
				}
			}
			if n.ShouldLog() {
				log.Warningf(ctx, "failed to perform initial scan: %v", err)
			}
		} else /* err == nil */ {
			if f.onInitialScanDone != nil {
				f.onInitialScanDone(ctx)
			}
			break
		}
	}

	return false
}

func (f *RangeFeed) getSpansToScan(ctx context.Context) func() []roachpb.Span {
	retryAll := func() []roachpb.Span {
		return f.spans
	}

	if f.retryBehavior == ScanRetryAll {
		return retryAll
	}

	// We want to retry remaining spans.
	// Maintain a frontier in order to keep track of which spans still need to be scanned.
	frontier, err := span.MakeFrontier(f.spans...)
	if err != nil {
		// Frontier construction shouldn't really fail. The spans have already
		// been validated by frontier constructed when starting rangefeed.
		// We don't really have a good mechanism to return an initialization error from here,
		// so, log it and fall back to retrying all spans.
		log.Errorf(ctx, "failed to build frontier for the initial scan; "+
			"falling back to retry all behavior: err=%v", err)
		return retryAll
	}

	var fm syncutil.Mutex
	userSpanDoneCallback := f.onSpanDone
	f.onSpanDone = func(ctx context.Context, sp roachpb.Span) error {
		if userSpanDoneCallback != nil {
			if err := userSpanDoneCallback(ctx, sp); err != nil {
				return err
			}
		}
		fm.Lock()
		defer fm.Unlock()
		_, err := frontier.Forward(sp, f.initialTimestamp)
		return err
	}

	isRetry := false
	var retrySpans []roachpb.Span
	return func() []roachpb.Span {
		if !isRetry {
			isRetry = true
			return f.spans
		}

		retrySpans = retrySpans[:0]
		frontier.Entries(func(sp roachpb.Span, ts hlc.Timestamp) (done span.OpResult) {
			if ts.IsEmpty() {
				retrySpans = append(retrySpans, sp)
			}
			return span.ContinueMatch
		})
		return retrySpans
	}

}
