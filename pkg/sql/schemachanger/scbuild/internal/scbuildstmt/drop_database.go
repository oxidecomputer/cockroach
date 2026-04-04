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

package scbuildstmt

import (
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/privilege"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/catid"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// DropDatabase implements DROP DATABASE.
func DropDatabase(b BuildCtx, n *tree.DropDatabase) {
	elts := b.ResolveDatabase(n.Name, ResolveParams{
		IsExistenceOptional: n.IfExists,
		RequiredPrivilege:   privilege.DROP,
	})
	_, _, db := scpb.FindDatabase(elts)
	if db == nil {
		return
	}
	if string(n.Name) == b.SessionData().Database && b.SessionData().SafeUpdates {
		panic(pgerror.DangerousStatementf("DROP DATABASE on current database"))
	}
	b.IncrementSchemaChangeDropCounter("database")
	// Perform explicit or implicit DROP DATABASE CASCADE.
	if n.DropBehavior == tree.DropCascade || (n.DropBehavior == tree.DropDefault && !b.SessionData().SafeUpdates) {
		dropCascadeDescriptor(b, db.DatabaseID)
		return
	}
	// Otherwise, perform DROP DATABASE RESTRICT.
	if !dropRestrictDescriptor(b, db.DatabaseID) {
		return
	}
	// Implicitly DROP RESTRICT the public schema as well.
	var publicSchemaID catid.DescID
	b.BackReferences(db.DatabaseID).ForEachElementStatus(func(_ scpb.Status, _ scpb.TargetStatus, e scpb.Element) {
		switch t := e.(type) {
		case *scpb.Schema:
			if t.IsPublic {
				publicSchemaID = t.SchemaID
			}
		}
	})
	dropRestrictDescriptor(b, publicSchemaID)
	dbBackrefs := undroppedBackrefs(b, db.DatabaseID)
	publicSchemaBackrefs := undroppedBackrefs(b, publicSchemaID)
	if dbBackrefs.IsEmpty() && publicSchemaBackrefs.IsEmpty() {
		return
	}
	// Block DROP if cascade is not set.
	if n.DropBehavior == tree.DropRestrict {
		panic(pgerror.Newf(pgcode.DependentObjectsStillExist,
			"database %q is not empty and RESTRICT was specified", simpleName(b, db.DatabaseID)))
	}
	panic(pgerror.DangerousStatementf(
		"DROP DATABASE on non-empty database without explicit CASCADE"))
}
