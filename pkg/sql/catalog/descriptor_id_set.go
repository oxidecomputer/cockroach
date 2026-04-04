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

package catalog

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/util"
)

// DescriptorIDSet efficiently stores an unordered set of descriptor ids.
type DescriptorIDSet struct {
	set util.FastIntSet
}

// MakeDescriptorIDSet returns a set initialized with the given values.
func MakeDescriptorIDSet(ids ...descpb.ID) DescriptorIDSet {
	s := DescriptorIDSet{}
	for _, id := range ids {
		s.Add(id)
	}
	return s
}

// Suppress the linter.
var _ = MakeDescriptorIDSet

// Add adds an id to the set. No-op if the id is already in the set.
func (d *DescriptorIDSet) Add(id descpb.ID) {
	d.set.Add(int(id))
}

// Len returns the number of the ids in the set.
func (d DescriptorIDSet) Len() int {
	return d.set.Len()
}

// Contains returns true if the set contains the column.
func (d DescriptorIDSet) Contains(id descpb.ID) bool {
	return d.set.Contains(int(id))
}

// ForEach calls a function for each column in the set (in increasing order).
func (d DescriptorIDSet) ForEach(f func(id descpb.ID)) {
	d.set.ForEach(func(i int) { f(descpb.ID(i)) })
}

// Empty returns true if the set is empty.
func (d DescriptorIDSet) Empty() bool { return d.set.Empty() }

// Ordered returns a slice with all the descpb.IDs in the set, in
// increasing order.
func (d DescriptorIDSet) Ordered() []descpb.ID {
	if d.Empty() {
		return nil
	}
	result := make([]descpb.ID, 0, d.Len())
	d.ForEach(func(i descpb.ID) {
		result = append(result, i)
	})
	return result
}

// Remove removes the ID from the set.
func (d *DescriptorIDSet) Remove(id descpb.ID) {
	d.set.Remove(int(id))
}
