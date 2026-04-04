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

package schemadesc

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catconstants"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
)

// GetVirtualSchemaByID returns a virtual schema with a given ID if it exists.
//
// The returned descriptor carries only a basic functionality, requiring the
// caller to check the SchemaKind to determine how to use the descriptor. The
// returned descriptor is not mapped to a database; every database has all of
// the same virtual schemas and the ParentID on the returned descriptor will be
// descpb.InvalidID.
func GetVirtualSchemaByID(id descpb.ID) (catalog.SchemaDescriptor, bool) {
	sc, ok := virtualSchemasByID[id]
	return sc, ok
}

var virtualSchemasByID = func() map[descpb.ID]catalog.SchemaDescriptor {
	m := make(map[descpb.ID]catalog.SchemaDescriptor, len(catconstants.StaticSchemaIDMap))
	for id, name := range catconstants.StaticSchemaIDMap {
		id := descpb.ID(id)
		sc := virtual{
			synthetic: synthetic{virtualBase{}},
			id:        id,
			name:      name,
		}
		m[id] = sc
	}
	return m
}()

// virtual represents the virtual schemas which are part of every database.
// See the commentary on GetVirtualSchemaByID.
type virtual struct {
	synthetic
	id   descpb.ID
	name string
}

var _ catalog.SchemaDescriptor = virtual{}

func (p virtual) GetID() descpb.ID       { return p.id }
func (p virtual) GetName() string        { return p.name }
func (p virtual) GetParentID() descpb.ID { return descpb.InvalidID }
func (p virtual) GetPrivileges() *catpb.PrivilegeDescriptor {
	return catpb.NewVirtualSchemaPrivilegeDescriptor()
}

type virtualBase struct{}

var _ syntheticBase = virtualBase{}

func (v virtualBase) kindName() string                 { return "virtual" }
func (v virtualBase) kind() catalog.ResolvedSchemaKind { return catalog.SchemaVirtual }
