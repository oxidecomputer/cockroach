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

package kvcoord

import (
	"sort"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// asSortedSlice returns the set data in sorted order.
//
// Too inefficient for production.
func (s *condensableSpanSet) asSortedSlice() []roachpb.Span {
	set := s.asSlice()
	cpy := make(roachpb.Spans, len(set))
	copy(cpy, set)
	sort.Sort(cpy)
	return cpy
}

// TestingSenderConcurrencyLimit exports the cluster setting for testing
// purposes.
var TestingSenderConcurrencyLimit = senderConcurrencyLimit

// TestingGetLockFootprint returns the internal lock footprint for testing
// purposes.
func (tc *TxnCoordSender) TestingGetLockFootprint(mergeAndSort bool) []roachpb.Span {
	if mergeAndSort {
		tc.interceptorAlloc.txnPipeliner.lockFootprint.mergeAndSort()
	}
	return tc.interceptorAlloc.txnPipeliner.lockFootprint.asSlice()
}

// TestingGetRefreshFootprint returns the internal refresh footprint for testing
// purposes.
func (tc *TxnCoordSender) TestingGetRefreshFootprint() []roachpb.Span {
	return tc.interceptorAlloc.txnSpanRefresher.refreshFootprint.asSlice()
}

// TestingSetLinearizable allows tests to enable linearizable behavior.
func (tcf *TxnCoordSenderFactory) TestingSetLinearizable(linearizable bool) {
	tcf.linearizable = linearizable
}

// TestingSetMetrics allows tests to override the factory's metrics struct.
func (tcf *TxnCoordSenderFactory) TestingSetMetrics(metrics TxnMetrics) {
	tcf.metrics = metrics
}
