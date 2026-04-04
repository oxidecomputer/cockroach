// Copyright 2018 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/sql/privilege"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scexec"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
)

type commentOnDatabaseNode struct {
	n               *tree.CommentOnDatabase
	dbDesc          catalog.DatabaseDescriptor
	metadataUpdater scexec.DescriptorMetadataUpdater
}

// CommentOnDatabase add comment on a database.
// Privileges: CREATE on database.
//   notes: postgres requires CREATE on the database.
func (p *planner) CommentOnDatabase(
	ctx context.Context, n *tree.CommentOnDatabase,
) (planNode, error) {
	if err := checkSchemaChangeEnabled(
		ctx,
		p.ExecCfg(),
		"COMMENT ON DATABASE",
	); err != nil {
		return nil, err
	}

	dbDesc, err := p.Descriptors().GetImmutableDatabaseByName(ctx, p.txn,
		string(n.Name), tree.DatabaseLookupFlags{Required: true})
	if err != nil {
		return nil, err
	}
	if err := p.CheckPrivilege(ctx, dbDesc, privilege.CREATE); err != nil {
		return nil, err
	}

	return &commentOnDatabaseNode{n: n,
		dbDesc: dbDesc,
		metadataUpdater: p.execCfg.DescMetadaUpdaterFactory.NewMetadataUpdater(
			ctx,
			p.txn,
			p.SessionData(),
		),
	}, nil
}

func (n *commentOnDatabaseNode) startExec(params runParams) error {
	if n.n.Comment != nil {
		err := n.metadataUpdater.UpsertDescriptorComment(
			int64(n.dbDesc.GetID()), 0, keys.DatabaseCommentType, *n.n.Comment)
		if err != nil {
			return err
		}
	} else {
		err := n.metadataUpdater.DeleteDescriptorComment(
			int64(n.dbDesc.GetID()), 0, keys.DatabaseCommentType)
		if err != nil {
			return err
		}
	}

	dbComment := ""
	if n.n.Comment != nil {
		dbComment = *n.n.Comment
	}
	return params.p.logEvent(params.ctx,
		n.dbDesc.GetID(),
		&eventpb.CommentOnDatabase{
			DatabaseName: n.n.Name.String(),
			Comment:      dbComment,
			NullComment:  n.n.Comment == nil,
		})
}

func (n *commentOnDatabaseNode) Next(runParams) (bool, error) { return false, nil }
func (n *commentOnDatabaseNode) Values() tree.Datums          { return tree.Datums{} }
func (n *commentOnDatabaseNode) Close(context.Context)        {}
