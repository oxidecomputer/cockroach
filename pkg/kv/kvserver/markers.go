// Copyright 2020 The Cockroach Authors.
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

package kvserver

import (
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/errors"
)

// NB: don't change the string here; this will cause cross-version issues
// since this singleton is used as a marker.
var errMarkSnapshotError = errors.New("snapshot failed")

// isSnapshotError returns true iff the error indicates that a snapshot failed.
func isSnapshotError(err error) bool {
	return errors.Is(err, errMarkSnapshotError)
}

// NB: don't change the string here; this will cause cross-version issues
// since this singleton is used as a marker.
var errMarkCanRetryReplicationChangeWithUpdatedDesc = errors.New("should retry with updated descriptor")

// IsRetriableReplicationChangeError detects whether an error (which is
// assumed to have been emitted by the KV layer during a replication change
// operation) is likely to succeed when retried with an updated descriptor.
func IsRetriableReplicationChangeError(err error) bool {
	return errors.Is(err, errMarkCanRetryReplicationChangeWithUpdatedDesc) || isSnapshotError(err)
}

const (
	descChangedErrorFmt              = "descriptor changed: [expected] %s != [actual] %s"
	descChangedRangeReplacedErrorFmt = "descriptor changed: [expected] %s != [actual] %s  (range replaced)"
	descChangedRangeSubsumedErrorFmt = "descriptor changed: [expected] %s != [actual] nil (range subsumed)"
)

func newDescChangedError(desc, actualDesc *roachpb.RangeDescriptor) error {
	var err error
	if actualDesc == nil {
		err = errors.Newf(descChangedRangeSubsumedErrorFmt, desc)
	} else if desc.RangeID != actualDesc.RangeID {
		err = errors.Newf(descChangedRangeReplacedErrorFmt, desc, actualDesc)
	} else {
		err = errors.Newf(descChangedErrorFmt, desc, actualDesc)
	}
	return errors.Mark(err, errMarkCanRetryReplicationChangeWithUpdatedDesc)
}

func wrapDescChangedError(err error, desc, actualDesc *roachpb.RangeDescriptor) error {
	var wrapped error
	if actualDesc == nil {
		wrapped = errors.Wrapf(err, descChangedRangeSubsumedErrorFmt, desc)
	} else if desc.RangeID != actualDesc.RangeID {
		wrapped = errors.Wrapf(err, descChangedRangeReplacedErrorFmt, desc, actualDesc)
	} else {
		wrapped = errors.Wrapf(err, descChangedErrorFmt, desc, actualDesc)
	}
	return errors.Mark(wrapped, errMarkCanRetryReplicationChangeWithUpdatedDesc)
}

// NB: don't change the string here; this will cause cross-version issues
// since this singleton is used as a marker.
var errMarkInvalidReplicationChange = errors.New("invalid replication change")

// IsIllegalReplicationChangeError detects whether an error (assumed to have been emitted
// by a replication change) indicates that the replication change is illegal, meaning
// that it's unlikely to be handled through a retry. Examples of this are attempts to add
// a store that is already a member of the supplied descriptor. A non-example is a change
// detected in the descriptor at the KV layer relative to that supplied as input to the
// replication change, which would likely benefit from a retry.
func IsIllegalReplicationChangeError(err error) bool {
	return errors.Is(err, errMarkInvalidReplicationChange)
}

var errMarkReplicationChangeInProgress = errors.New("replication change in progress")

// IsReplicationChangeInProgressError detects whether an error (assumed to have
// been emitted a replication change) indicates that the replication change
// failed because another replication change was in progress on the range.
func IsReplicationChangeInProgressError(err error) bool {
	return errors.Is(err, errMarkReplicationChangeInProgress)
}
