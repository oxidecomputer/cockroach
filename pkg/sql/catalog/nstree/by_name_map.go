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

type byNameMap struct {
	t *btree.BTree
}

func (t byNameMap) upsert(d catalog.NameEntry) (replaced catalog.NameEntry) {
	replaced, _ = upsert(t.t, makeByNameItem(d).get()).(catalog.NameEntry)
	return replaced
}

func (t byNameMap) getByName(parentID, parentSchemaID descpb.ID, name string) catalog.NameEntry {
	got, _ := get(t.t, byNameItem{
		parentID:       parentID,
		parentSchemaID: parentSchemaID,
		name:           name,
	}.get()).(catalog.NameEntry)
	return got
}

func (t byNameMap) delete(d catalog.NameKey) (removed catalog.NameEntry) {
	removed, _ = delete(t.t, makeByNameItem(d).get()).(catalog.NameEntry)
	return removed
}

func (t byNameMap) clear() {
	clear(t.t)
}

func (t byNameMap) ascend(f EntryIterator) error {
	return ascend(t.t, func(k interface{}) error {
		return f(k.(catalog.NameEntry))
	})
}
