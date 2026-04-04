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

package nstree

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/google/btree"
)

// Map is a lookup structure for descriptors. It is used to provide
// indexed access to a set of entries either by name or by ID. The
// entries' properties are indexed; they must not change or else the
// index will be corrupted. Safe for use without initialization. Calling
// Clear will return memory to a sync.Pool.
type Map struct {
	byID   byIDMap
	byName byNameMap
}

// EntryIterator is used to iterate namespace entries.
// If an error is returned, iteration is stopped and will be propagated
// up the stack. If the error is iterutil.StopIteration, iteration will
// stop but no error will be returned.
type EntryIterator func(entry catalog.NameEntry) error

// Upsert adds the descriptor to the tree. If any descriptor exists in the
// tree with the same name or id, it will be removed.
func (dt *Map) Upsert(d catalog.NameEntry) {
	dt.maybeInitialize()
	if replaced := dt.byName.upsert(d); replaced != nil {
		dt.byID.delete(replaced.GetID())
	}
	if replaced := dt.byID.upsert(d); replaced != nil {
		dt.byName.delete(replaced)
	}
}

// Remove removes the descriptor with the given ID from the tree and
// returns it if it exists.
func (dt *Map) Remove(id descpb.ID) catalog.NameEntry {
	dt.maybeInitialize()
	if d := dt.byID.delete(id); d != nil {
		dt.byName.delete(d)
		return d
	}
	return nil
}

// GetByID gets a descriptor from the tree by id.
func (dt *Map) GetByID(id descpb.ID) catalog.NameEntry {
	if !dt.initialized() {
		return nil
	}
	return dt.byID.get(id)
}

// GetByName gets a descriptor from the tree by name.
func (dt *Map) GetByName(parentID, parentSchemaID descpb.ID, name string) catalog.NameEntry {
	if !dt.initialized() {
		return nil
	}
	return dt.byName.getByName(parentID, parentSchemaID, name)
}

// Clear removes all entries, returning any held memory to the sync.Pool.
func (dt *Map) Clear() {
	if !dt.initialized() {
		return
	}
	dt.byID.clear()
	dt.byName.clear()
	btreeSyncPool.Put(dt.byName.t)
	btreeSyncPool.Put(dt.byID.t)
	*dt = Map{}
}

// IterateByID iterates the descriptors by ID, ascending.
func (dt *Map) IterateByID(f EntryIterator) error {
	if !dt.initialized() {
		return nil
	}
	return dt.byID.ascend(f)
}

// IterateByName iterates the descriptors by name, ascending.
func (dt *Map) IterateByName(f EntryIterator) error {
	if !dt.initialized() {
		return nil
	}
	return dt.byName.ascend(f)
}

// Len returns the number of descriptors in the tree.
func (dt *Map) Len() int {
	if !dt.initialized() {
		return 0
	}
	return dt.byID.len()
}

func (dt Map) initialized() bool {
	return dt != (Map{})
}

func (dt *Map) maybeInitialize() {
	if dt.initialized() {
		return
	}
	*dt = Map{
		byName: byNameMap{t: btreeSyncPool.Get().(*btree.BTree)},
		byID:   byIDMap{t: btreeSyncPool.Get().(*btree.BTree)},
	}
}
