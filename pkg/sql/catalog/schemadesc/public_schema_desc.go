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
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// GetPublicSchema returns a synthetic public schema which is
// part of every database. The public schema's implementation is a vestige
// of a time when there were no user-defined schemas. The public schema is
// interchangeable with the database itself in terms of privileges.
//
// The returned descriptor carries only a basic functionality, requiring the
// caller to check the SchemaKind to determine how to use the descriptor. The
// returned descriptor is not mapped to a database; every database has all of
// the same virtual schemas and the ParentID on the returned descriptor will be
// descpb.InvalidID.
// This is deprecated and should not be used except for certain edge cases.
// This will be removed in 22.2 completely.
func GetPublicSchema() catalog.SchemaDescriptor {
	return publicDesc
}

type public struct {
	synthetic
}

var _ catalog.SchemaDescriptor = public{}

func (p public) GetParentID() descpb.ID { return descpb.InvalidID }
func (p public) GetID() descpb.ID       { return keys.PublicSchemaID }
func (p public) GetName() string        { return tree.PublicSchema }
func (p public) GetPrivileges() *catpb.PrivilegeDescriptor {
	return catpb.NewPublicSchemaPrivilegeDescriptor()
}

type publicBase struct{}

func (p publicBase) kindName() string                 { return "public" }
func (p publicBase) kind() catalog.ResolvedSchemaKind { return catalog.SchemaPublic }

// publicDesc is a singleton returned by GetPublicSchema.
var publicDesc catalog.SchemaDescriptor = public{
	synthetic{publicBase{}},
}
