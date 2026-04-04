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

package kvserver

import (
	"bytes"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/redact"
)

// ReplicaSnapshotDiff is a part of a []ReplicaSnapshotDiff which represents a diff between
// two replica snapshots. For now it's only a diff between their KV pairs.
type ReplicaSnapshotDiff struct {
	// LeaseHolder is set to true of this kv pair is only present on the lease
	// holder.
	LeaseHolder bool
	Key         roachpb.Key
	Timestamp   hlc.Timestamp
	Value       []byte
}

// ReplicaSnapshotDiffSlice groups multiple ReplicaSnapshotDiff records and
// exposes a formatting helper.
type ReplicaSnapshotDiffSlice []ReplicaSnapshotDiff

// SafeFormat implements redact.SafeFormatter.
func (rsds ReplicaSnapshotDiffSlice) SafeFormat(buf redact.SafePrinter, _ rune) {
	buf.Printf("--- leaseholder\n+++ follower\n")
	for _, d := range rsds {
		prefix := redact.SafeString("+")
		if d.LeaseHolder {
			// Lease holder (RHS) has something follower (LHS) does not have.
			prefix = redact.SafeString("-")
		}
		const format = `%s%s %s
%s    ts:%s
%s    value:%s
%s    raw mvcc_key/value: %x %x
`
		mvccKey := storage.MVCCKey{Key: d.Key, Timestamp: d.Timestamp}
		buf.Printf(format,
			prefix, d.Timestamp, d.Key,
			prefix, d.Timestamp.GoTime(),
			prefix, SprintMVCCKeyValue(storage.MVCCKeyValue{Key: mvccKey, Value: d.Value}, false /* printKey */),
			prefix, storage.EncodeMVCCKey(mvccKey), d.Value)
	}
}

func (rsds ReplicaSnapshotDiffSlice) String() string {
	return redact.StringWithoutMarkers(rsds)
}

// diffs the two kv dumps between the lease holder and the replica.
func diffRange(l, r *roachpb.RaftSnapshotData) ReplicaSnapshotDiffSlice {
	if l == nil || r == nil {
		return nil
	}
	var diff []ReplicaSnapshotDiff
	i, j := 0, 0
	for {
		var e, v roachpb.RaftSnapshotData_KeyValue
		if i < len(l.KV) {
			e = l.KV[i]
		}
		if j < len(r.KV) {
			v = r.KV[j]
		}

		addLeaseHolder := func() {
			diff = append(diff, ReplicaSnapshotDiff{LeaseHolder: true, Key: e.Key, Timestamp: e.Timestamp, Value: e.Value})
			i++
		}
		addReplica := func() {
			diff = append(diff, ReplicaSnapshotDiff{LeaseHolder: false, Key: v.Key, Timestamp: v.Timestamp, Value: v.Value})
			j++
		}

		// Compare keys.
		var comp int
		// Check if it has finished traversing over all the lease holder keys.
		if e.Key == nil {
			if v.Key == nil {
				// Done traversing over all the replica keys. Done!
				break
			} else {
				comp = 1
			}
		} else {
			// Check if it has finished traversing over all the replica keys.
			if v.Key == nil {
				comp = -1
			} else {
				// Both lease holder and replica keys exist. Compare them.
				comp = bytes.Compare(e.Key, v.Key)
			}
		}
		switch comp {
		case -1:
			addLeaseHolder()

		case 0:
			// Timestamp sorting is weird. Timestamp{} sorts first, the
			// remainder sort in descending order. See storage/engine/doc.go.
			if !e.Timestamp.EqOrdering(v.Timestamp) {
				if e.Timestamp.IsEmpty() {
					addLeaseHolder()
				} else if v.Timestamp.IsEmpty() {
					addReplica()
				} else if v.Timestamp.Less(e.Timestamp) {
					addLeaseHolder()
				} else {
					addReplica()
				}
			} else if !bytes.Equal(e.Value, v.Value) {
				addLeaseHolder()
				addReplica()
			} else {
				// No diff; skip.
				i++
				j++
			}

		case 1:
			addReplica()

		}
	}
	return diff
}
