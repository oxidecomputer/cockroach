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

package sql

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/privilege"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scexec"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

type commentOnSchemaNode struct {
	n               *tree.CommentOnSchema
	schemaDesc      catalog.SchemaDescriptor
	metadataUpdater scexec.DescriptorMetadataUpdater
}

// CommentOnSchema add comment on a schema.
// Privileges: CREATE on scheme.
//   notes: postgres requires CREATE on the scheme.
func (p *planner) CommentOnSchema(ctx context.Context, n *tree.CommentOnSchema) (planNode, error) {
	if err := checkSchemaChangeEnabled(
		ctx,
		p.ExecCfg(),
		"COMMENT ON SCHEMA",
	); err != nil {
		return nil, err
	}

	// Users can't create a schema without being connected to a DB.
	dbName := p.CurrentDatabase()
	if dbName == "" {
		return nil, pgerror.New(pgcode.UndefinedDatabase,
			"cannot comment schema without being connected to a database")
	}

	db, err := p.Descriptors().GetImmutableDatabaseByName(ctx, p.txn,
		dbName, tree.DatabaseLookupFlags{Required: true})
	if err != nil {
		return nil, err
	}

	schemaDesc, err := p.Descriptors().GetImmutableSchemaByID(ctx, p.txn,
		db.GetSchemaID(string(n.Name)), tree.SchemaLookupFlags{Required: true})
	if err != nil {
		return nil, err
	}

	if err := p.CheckPrivilege(ctx, db, privilege.CREATE); err != nil {
		return nil, err
	}

	return &commentOnSchemaNode{
		n:          n,
		schemaDesc: schemaDesc,
		metadataUpdater: p.execCfg.DescMetadaUpdaterFactory.NewMetadataUpdater(
			ctx,
			p.txn,
			p.SessionData(),
		),
	}, nil
}

func (n *commentOnSchemaNode) startExec(params runParams) error {
	if n.n.Comment != nil {
		err := n.metadataUpdater.UpsertDescriptorComment(
			int64(n.schemaDesc.GetID()), 0, keys.SchemaCommentType, *n.n.Comment)
		if err != nil {
			return err
		}
	} else {
		err := n.metadataUpdater.DeleteDescriptorComment(
			int64(n.schemaDesc.GetID()), 0, keys.SchemaCommentType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *commentOnSchemaNode) Next(runParams) (bool, error) { return false, nil }
func (n *commentOnSchemaNode) Values() tree.Datums          { return tree.Datums{} }
func (n *commentOnSchemaNode) Close(context.Context)        {}
