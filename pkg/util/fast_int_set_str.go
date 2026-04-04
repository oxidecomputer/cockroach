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

package util

import (
	"bytes"
	"fmt"
)

// String returns a list representation of elements. Sequential runs of positive
// numbers are shown as ranges. For example, for the set {0, 1, 2, 5, 6, 10},
// the output is "(0-2,5,6,10)".
func (s FastIntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('(')
	appendRange := func(start, end int) {
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}
		if start == end {
			fmt.Fprintf(&buf, "%d", start)
		} else if start+1 == end {
			fmt.Fprintf(&buf, "%d,%d", start, end)
		} else {
			fmt.Fprintf(&buf, "%d-%d", start, end)
		}
	}
	rangeStart, rangeEnd := -1, -1
	s.ForEach(func(i int) {
		if i < 0 {
			appendRange(i, i)
			return
		}
		if rangeStart != -1 && rangeEnd == i-1 {
			rangeEnd = i
		} else {
			if rangeStart != -1 {
				appendRange(rangeStart, rangeEnd)
			}
			rangeStart, rangeEnd = i, i
		}
	})
	if rangeStart != -1 {
		appendRange(rangeStart, rangeEnd)
	}
	buf.WriteByte(')')
	return buf.String()
}
