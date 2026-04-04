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

package schemaexpr_test

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/tabledesc"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// testCol includes the information needed to create a column descriptor for
// testing purposes.
type testCol struct {
	name string
	typ  *types.T
}

// testTableDesc is a helper functions for creating table descriptors in a
// less verbose way.
func testTableDesc(
	name string, columns []testCol, mutationColumns []testCol,
) catalog.MutableTableDescriptor {
	cols := make([]descpb.ColumnDescriptor, len(columns))
	for i := range columns {
		cols[i] = descpb.ColumnDescriptor{
			Name: columns[i].name,
			Type: columns[i].typ,
			// Column IDs start at 1 to mimic "real" table descriptors.
			ID: descpb.ColumnID(i + 1),
		}
	}

	muts := make([]descpb.DescriptorMutation, len(mutationColumns))
	for i := range mutationColumns {
		muts[i] = descpb.DescriptorMutation{
			Descriptor_: &descpb.DescriptorMutation_Column{
				Column: &descpb.ColumnDescriptor{
					Name: mutationColumns[i].name,
					Type: mutationColumns[i].typ,
					ID:   descpb.ColumnID(len(columns) + i + 1),
				},
			},
			Direction: descpb.DescriptorMutation_ADD,
		}
	}
	return tabledesc.NewBuilder(&descpb.TableDescriptor{
		Name:      name,
		ID:        1,
		Columns:   cols,
		Mutations: muts,
	}).BuildCreatedMutableTable()
}
