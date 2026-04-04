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

type commentOnTableNode struct {
	n               *tree.CommentOnTable
	tableDesc       catalog.TableDescriptor
	metadataUpdater scexec.DescriptorMetadataUpdater
}

// CommentOnTable add comment on a table.
// Privileges: CREATE on table.
//   notes: postgres requires CREATE on the table.
//          mysql requires ALTER, CREATE, INSERT on the table.
func (p *planner) CommentOnTable(ctx context.Context, n *tree.CommentOnTable) (planNode, error) {
	if err := checkSchemaChangeEnabled(
		ctx,
		p.ExecCfg(),
		"COMMENT ON TABLE",
	); err != nil {
		return nil, err
	}

	tableDesc, err := p.ResolveUncachedTableDescriptorEx(ctx, n.Table, true, tree.ResolveRequireTableDesc)
	if err != nil {
		return nil, err
	}

	if err := p.CheckPrivilege(ctx, tableDesc, privilege.CREATE); err != nil {
		return nil, err
	}

	return &commentOnTableNode{
		n:         n,
		tableDesc: tableDesc,
		metadataUpdater: p.execCfg.DescMetadaUpdaterFactory.NewMetadataUpdater(
			ctx,
			p.txn,
			p.SessionData(),
		),
	}, nil
}

func (n *commentOnTableNode) startExec(params runParams) error {
	if n.n.Comment != nil {
		err := n.metadataUpdater.UpsertDescriptorComment(
			int64(n.tableDesc.GetID()), 0, keys.TableCommentType, *n.n.Comment)
		if err != nil {
			return err
		}
	} else {
		err := n.metadataUpdater.DeleteDescriptorComment(
			int64(n.tableDesc.GetID()), 0, keys.TableCommentType)
		if err != nil {
			return err
		}
	}

	comment := ""
	if n.n.Comment != nil {
		comment = *n.n.Comment
	}
	return params.p.logEvent(params.ctx,
		n.tableDesc.GetID(),
		&eventpb.CommentOnTable{
			TableName:   params.p.ResolvedName(n.n.Table).FQString(),
			Comment:     comment,
			NullComment: n.n.Comment == nil,
		})
}

func (n *commentOnTableNode) Next(runParams) (bool, error) { return false, nil }
func (n *commentOnTableNode) Values() tree.Datums          { return tree.Datums{} }
func (n *commentOnTableNode) Close(context.Context)        {}
