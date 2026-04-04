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

package sql

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/tabledesc"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
)

type alterTableOwnerNode struct {
	owner  security.SQLUsername
	desc   *tabledesc.Mutable
	n      *tree.AlterTableOwner
	prefix catalog.ResolvedObjectPrefix
}

// AlterTableOwner sets the owner for a table, view, or sequence.
func (p *planner) AlterTableOwner(ctx context.Context, n *tree.AlterTableOwner) (planNode, error) {
	if err := checkSchemaChangeEnabled(
		ctx,
		p.ExecCfg(),
		"ALTER TABLE/VIEW/SEQUENCE OWNER",
	); err != nil {
		return nil, err
	}

	tn := n.Name.ToTableName()

	// ALTER TABLE [table] OWNER TO command applies to VIEWS and SEQUENCES
	// so we must resolve any table kind.
	requiredTableKind := tree.ResolveAnyTableKind
	if n.IsView {
		requiredTableKind = tree.ResolveRequireViewDesc
	} else if n.IsSequence {
		requiredTableKind = tree.ResolveRequireSequenceDesc
	}
	prefix, tableDesc, err := p.ResolveMutableTableDescriptor(
		ctx, &tn, !n.IfExists, requiredTableKind)
	if err != nil {
		return nil, err
	}
	if tableDesc == nil {
		// Noop.
		return newZeroNode(nil /* columns */), nil
	}

	if err := checkViewMatchesMaterialized(tableDesc, n.IsView, n.IsMaterialized); err != nil {
		return nil, err
	}

	owner, err := n.Owner.ToSQLUsername(p.SessionData(), security.UsernameValidation)
	if err != nil {
		return nil, err
	}
	return &alterTableOwnerNode{
		owner:  owner,
		desc:   tableDesc,
		n:      n,
		prefix: prefix,
	}, nil
}

func (n *alterTableOwnerNode) startExec(params runParams) error {
	telemetry.Inc(n.n.TelemetryCounter())
	ctx := params.ctx
	p := params.p
	tableDesc := n.desc
	newOwner := n.owner
	oldOwner := n.desc.GetPrivileges().Owner()

	if err := p.checkCanAlterToNewOwner(ctx, tableDesc, newOwner); err != nil {
		return err
	}

	// Ensure the new owner has CREATE privilege on the table's schema.
	if err := p.canCreateOnSchema(
		ctx, tableDesc.GetParentSchemaID(), tableDesc.ParentID, newOwner, checkPublicSchema); err != nil {
		return err
	}

	tbNameWithSchema := tree.MakeTableNameWithSchema(
		tree.Name(n.prefix.Database.GetName()),
		tree.Name(n.prefix.Schema.GetName()),
		tree.Name(tableDesc.GetName()),
	)

	if err := p.setNewTableOwner(ctx, tableDesc, tbNameWithSchema, newOwner); err != nil {
		return err
	}

	// If the owner we want to set to is the current owner, do a no-op.
	if newOwner == oldOwner {
		return nil
	}

	if err := p.writeSchemaChange(
		ctx, tableDesc, descpb.InvalidMutationID, tree.AsStringWithFQNames(n.n, params.Ann()),
	); err != nil {
		return err
	}

	return nil
}

// setNewTableOwner handles setting a new table owner.
// Called in ALTER SCHEMA and REASSIGN OWNED BY.
func (p *planner) setNewTableOwner(
	ctx context.Context,
	desc *tabledesc.Mutable,
	tbNameWithSchema tree.TableName,
	newOwner security.SQLUsername,
) error {
	privs := desc.GetPrivileges()
	privs.SetOwner(newOwner)

	return p.logEvent(ctx,
		desc.ID,
		&eventpb.AlterTableOwner{
			TableName: tbNameWithSchema.FQString(),
			Owner:     newOwner.Normalized(),
		})
}

// ReadingOwnWrites implements the planNodeReadingOwnWrites interface.
// This is because SET SCHEMA performs multiple KV operations on descriptors
// and expects to see its own writes.
func (n *alterTableOwnerNode) ReadingOwnWrites() {}

func (n *alterTableOwnerNode) Next(runParams) (bool, error) { return false, nil }
func (n *alterTableOwnerNode) Values() tree.Datums          { return tree.Datums{} }
func (n *alterTableOwnerNode) Close(context.Context)        {}
