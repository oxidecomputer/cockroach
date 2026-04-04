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
	"sync"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/google/btree"
)

type byNameItem struct {
	parentID, parentSchemaID descpb.ID
	name                     string
	v                        interface{}
}

func makeByNameItem(d catalog.NameKey) byNameItem {
	return byNameItem{
		parentID:       d.GetParentID(),
		parentSchemaID: d.GetParentSchemaID(),
		name:           d.GetName(),
		v:              d,
	}
}

var _ btree.Item = (*byNameItem)(nil)

func (b *byNameItem) Less(thanItem btree.Item) bool {
	than := thanItem.(*byNameItem)
	if b.parentID != than.parentID {
		return b.parentID < than.parentID
	}
	if b.parentSchemaID != than.parentSchemaID {
		return b.parentSchemaID < than.parentSchemaID
	}
	return b.name < than.name
}

func (b *byNameItem) value() interface{} {
	return b.v
}

var byNameItemPool = sync.Pool{
	New: func() interface{} { return new(byNameItem) },
}

func (b byNameItem) get() *byNameItem {
	item := byNameItemPool.Get().(*byNameItem)
	*item = b
	return item
}

func (b *byNameItem) put() {
	*b = byNameItem{}
	byNameItemPool.Put(b)
}
