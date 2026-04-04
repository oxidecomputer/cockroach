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
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/privilege"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
)

type commentOnColumnNode struct {
	n         *tree.CommentOnColumn
	tableDesc catalog.TableDescriptor
}

// CommentOnColumn add comment on a column.
// Privileges: CREATE on table.
func (p *planner) CommentOnColumn(ctx context.Context, n *tree.CommentOnColumn) (planNode, error) {
	if err := checkSchemaChangeEnabled(
		ctx,
		p.ExecCfg(),
		"COMMENT ON COLUMN",
	); err != nil {
		return nil, err
	}

	var tableName tree.TableName
	if n.ColumnItem.TableName != nil {
		tableName = n.ColumnItem.TableName.ToTableName()
	}
	tableDesc, err := p.resolveUncachedTableDescriptor(ctx, &tableName, true, tree.ResolveRequireTableDesc)
	if err != nil {
		return nil, err
	}

	if err := p.CheckPrivilege(ctx, tableDesc, privilege.CREATE); err != nil {
		return nil, err
	}

	return &commentOnColumnNode{n: n, tableDesc: tableDesc}, nil
}

func (n *commentOnColumnNode) startExec(params runParams) error {
	col, err := n.tableDesc.FindColumnWithName(n.n.ColumnItem.ColumnName)
	if err != nil {
		return err
	}

	if n.n.Comment != nil {
		_, err := params.p.extendedEvalCtx.ExecCfg.InternalExecutor.ExecEx(
			params.ctx,
			"set-column-comment",
			params.p.Txn(),
			sessiondata.InternalExecutorOverride{User: security.RootUserName()},
			"UPSERT INTO system.comments VALUES ($1, $2, $3, $4)",
			keys.ColumnCommentType,
			n.tableDesc.GetID(),
			col.GetPGAttributeNum(),
			*n.n.Comment)
		if err != nil {
			return err
		}
	} else {
		_, err := params.p.extendedEvalCtx.ExecCfg.InternalExecutor.ExecEx(
			params.ctx,
			"delete-column-comment",
			params.p.Txn(),
			sessiondata.InternalExecutorOverride{User: security.RootUserName()},
			"DELETE FROM system.comments WHERE type=$1 AND object_id=$2 AND sub_id=$3",
			keys.ColumnCommentType,
			n.tableDesc.GetID(),
			col.GetPGAttributeNum())
		if err != nil {
			return err
		}
	}

	comment := ""
	if n.n.Comment != nil {
		comment = *n.n.Comment
	}

	tn, err := params.p.getQualifiedTableName(params.ctx, n.tableDesc)
	if err != nil {
		return err
	}

	return params.p.logEvent(params.ctx,
		n.tableDesc.GetID(),
		&eventpb.CommentOnColumn{
			TableName:   tn.FQString(),
			ColumnName:  string(n.n.ColumnItem.ColumnName),
			Comment:     comment,
			NullComment: n.n.Comment == nil,
		})
}

func (n *commentOnColumnNode) Next(runParams) (bool, error) { return false, nil }
func (n *commentOnColumnNode) Values() tree.Datums          { return tree.Datums{} }
func (n *commentOnColumnNode) Close(context.Context)        {}
