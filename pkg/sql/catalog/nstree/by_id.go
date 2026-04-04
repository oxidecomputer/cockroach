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

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/google/btree"
)

func (t byIDMap) len() int {
	return t.t.Len()
}

type byIDItem struct {
	id descpb.ID
	v  interface{}
}

func makeByIDItem(d interface{ GetID() descpb.ID }) byIDItem {
	return byIDItem{id: d.GetID(), v: d}
}

var _ btree.Item = (*byIDItem)(nil)

func (b *byIDItem) Less(thanItem btree.Item) bool {
	than := thanItem.(*byIDItem)
	return b.id < than.id
}

var byIDItemPool = sync.Pool{
	New: func() interface{} { return new(byIDItem) },
}

func (b byIDItem) get() *byIDItem {
	alloc := byIDItemPool.Get().(*byIDItem)
	*alloc = b
	return alloc
}

func (b *byIDItem) value() interface{} {
	return b.v
}

func (b *byIDItem) put() {
	*b = byIDItem{}
	byIDItemPool.Put(b)
}
