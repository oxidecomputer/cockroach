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
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
)

// NewTemporarySchema returns a temporary schema with a given name, id, and
// parent. Temporary schemas do not have a durable descriptor in the store;
// they only have a namespace entry to indicate their existence. Given that,
// a different kind of "synthetic" descriptor is used to indicate temporary
// schemas.
//
// The returned descriptor carries only a basic functionality, requiring the
// caller to check the SchemaKind to determine how to use the descriptor.
func NewTemporarySchema(name string, id descpb.ID, parentDB descpb.ID) catalog.SchemaDescriptor {
	return &temporary{
		synthetic: synthetic{temporaryBase{}},
		id:        id,
		name:      name,
		parentID:  parentDB,
	}
}

// temporary represents the synthetic temporary schema.
type temporary struct {
	synthetic
	id       descpb.ID
	name     string
	parentID descpb.ID
}

var _ catalog.SchemaDescriptor = temporary{}

func (p temporary) GetID() descpb.ID       { return p.id }
func (p temporary) GetName() string        { return p.name }
func (p temporary) GetParentID() descpb.ID { return p.parentID }
func (p temporary) GetPrivileges() *catpb.PrivilegeDescriptor {
	return catpb.NewTemporarySchemaPrivilegeDescriptor()
}

type temporaryBase struct{}

func (t temporaryBase) kindName() string                 { return "temporary" }
func (t temporaryBase) kind() catalog.ResolvedSchemaKind { return catalog.SchemaTemporary }

var _ syntheticBase = temporaryBase{}
