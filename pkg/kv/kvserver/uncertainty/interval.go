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

package uncertainty

import "github.com/cockroachdb/cockroach/pkg/util/hlc"

// Interval represents a transaction's uncertainty interval.
//
// The GlobalLimit is the transaction's initial upper bound on the interval,
// which is computed by its coordinator when the transaction is initiated as the
// current local HLC time plus the maximum clock offset between any two nodes in
// the system. Assuming this maximum clock skew bound holds, this represents an
// upper bound on the HLC time of any node in the system at the instant that the
// transaction began.
//
// The LocalLimit is an optional tighter bound established through HLC clock
// observations on individual nodes in the system. These clock observations
// ("observed timestamps") can only be used to lower the local limit of an
// uncertainty interval when a transaction evaluates on the node from which the
// observation was taken.
//
// The local limit can reduce the uncertainty interval applied to most values on
// a range. This can lead to values that would otherwise be considered uncertain
// by the original global limit to be considered "certainly concurrent", and
// thus not causally related, with the transaction due to observed timestamps.
//
// However, the local limit does not apply to all committed values on a range.
// Specifically, values with "synthetic timestamps" must use the interval's
// global limit for the purposes of uncertainty, because observed timestamps do
// not apply to values with synthetic timestamps.
//
// Uncertainty intervals also apply to non-transactional requests that require
// strong consistency (single-key linearizability). These requests defer their
// timestamp allocation to the leaseholder of their (single) range. They then
// establish an uncertainty interval and perform any uncertainty restarts on the
// server.
//
// NOTE: LocalLimit can be empty if no observed timestamp has been captured on
// the local node. However, if it is set, it must be <= GlobalLimit.
type Interval struct {
	GlobalLimit hlc.Timestamp
	LocalLimit  hlc.ClockTimestamp
}

// IsUncertain determines whether a value with the provided timestamp is
// uncertain to a reader with a ReadTimestamp below the value's and with
// the specified uncertainty interval.
func (in Interval) IsUncertain(valueTs hlc.Timestamp) bool {
	if !in.LocalLimit.IsEmpty() && !valueTs.Synthetic {
		return valueTs.LessEq(in.LocalLimit.ToTimestamp())
	}
	return valueTs.LessEq(in.GlobalLimit)
}
