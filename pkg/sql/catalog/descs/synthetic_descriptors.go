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

package descs

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/nstree"
)

type syntheticDescriptors struct {
	descs nstree.Map
}

func (sd *syntheticDescriptors) add(desc catalog.Descriptor) {
	if mut, ok := desc.(catalog.MutableDescriptor); ok {
		desc = mut.ImmutableCopy()
		sd.descs.Upsert(desc)
	} else {
		// Already an immutable object.
		sd.descs.Upsert(desc)
	}
}

func (sd *syntheticDescriptors) remove(id descpb.ID) {
	sd.descs.Remove(id)
}

func (sd *syntheticDescriptors) set(descs []catalog.Descriptor) {
	sd.descs.Clear()
	for _, desc := range descs {
		sd.add(desc)
	}
}

func (sd *syntheticDescriptors) reset() {
	sd.descs.Clear()
}

func (sd *syntheticDescriptors) getByName(
	dbID descpb.ID, schemaID descpb.ID, name string,
) (found bool, desc catalog.Descriptor) {
	if entry := sd.descs.GetByName(dbID, schemaID, name); entry != nil {
		return true, entry.(catalog.Descriptor)
	}
	return false, nil
}

func (sd *syntheticDescriptors) getByID(id descpb.ID) (found bool, desc catalog.Descriptor) {
	if entry := sd.descs.GetByID(id); entry != nil {
		return true, entry.(catalog.Descriptor)
	}
	return false, nil
}
