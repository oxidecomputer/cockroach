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

package rel

import "math/bits"

// ordinal is used to correlate attributes in a schema.
// It enables use of the ordinalSet.
type ordinal uint64

// ordinalSet represents A bitmask over ordinals.
// Note that it cannot contain attributes with ordinals greater than 64.
type ordinalSet uint64

// ForEach iterates the set of attributes.
func (m ordinalSet) forEach(f func(a ordinal) (wantMore bool)) {
	rem := m
	for rem > 0 {
		ord := ordinal(bits.TrailingZeros64(uint64(rem)))
		if !f(ord) {
			return
		}
		rem = rem.remove(ord)
	}
}

// remove returns the set constructed by removing ord from m.
func (m ordinalSet) remove(ord ordinal) ordinalSet {
	return m & ^(1 << ord)
}

// contains tests if m contains ord.
func (m ordinalSet) contains(ord ordinal) bool {
	return m&(1<<ord) != 0
}

// add returns the set constructed by adding ord to m.
func (m ordinalSet) add(ord ordinal) ordinalSet {
	return m | (1 << ord)
}

// without returns the set constructed by removing the members of other from m.
func (m ordinalSet) without(other ordinalSet) ordinalSet {
	return m & ^other
}

// intersection returns the set constructing with the intersection of m and other.
func (m ordinalSet) intersection(other ordinalSet) ordinalSet {
	return m & other
}

// union returns the set constructing with the union of m and other.
func (m ordinalSet) union(other ordinalSet) ordinalSet {
	return m | other
}

// len returns the number of ordinals in the set.
func (m ordinalSet) len() int {
	return bits.OnesCount64(uint64(m))
}
