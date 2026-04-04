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

package ctpb

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
)

// SeqNum identifies a ctpb.Update.
type SeqNum int64

// SafeValue implements the redact.SafeValue interface.
func (SeqNum) SafeValue() {}

// LAI is an int64 denoting a lease applied index with its own type to avoid
// mix-ups in positional arguments.
type LAI int64

// SafeValue implements the redact.SafeValue interface.
func (LAI) SafeValue() {}

func (m *Update) String() string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "Seq num: %d, sending node: n%d, snapshot: %t, size: %d bytes",
		m.SeqNum, m.NodeID, m.Snapshot, m.Size())
	sb.WriteString(", closed timestamps: ")
	now := timeutil.Now()
	for i, upd := range m.ClosedTimestamps {
		if i != 0 {
			sb.WriteString(", ")
		}
		ago := now.Sub(upd.ClosedTimestamp.GoTime()).Truncate(time.Millisecond)
		var agoMsg string
		if ago >= 0 {
			agoMsg = fmt.Sprintf("%s ago", ago)
		} else {
			agoMsg = fmt.Sprintf("%s in the future", -ago)
		}
		fmt.Fprintf(sb, "%s:%s (%s)", upd.Policy, upd.ClosedTimestamp, agoMsg)
	}
	sb.WriteRune('\n')

	fmt.Fprintf(sb, "Added or updated (%d ranges): (<range>:<LAI>) ", len(m.AddedOrUpdated))
	added := make([]Update_RangeUpdate, len(m.AddedOrUpdated))
	copy(added, m.AddedOrUpdated)
	sort.Slice(added, func(i, j int) bool {
		return added[i].RangeID < added[j].RangeID
	})
	for i, upd := range m.AddedOrUpdated {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(sb, "%d:%d", upd.RangeID, upd.LAI)
	}
	sb.WriteRune('\n')

	fmt.Fprintf(sb, "Removed (%d ranges): ", len(m.Removed))
	removed := make([]roachpb.RangeID, len(m.Removed))
	copy(removed, m.Removed)
	sort.Slice(removed, func(i, j int) bool {
		return removed[i] < removed[j]
	})
	for i, rid := range removed {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(sb, "r%d", rid)
	}
	sb.WriteRune('\n')
	return sb.String()
}
