// Copyright 2022 The Cockroach Authors.
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

package scmutationexec

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/screl"
	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
)

func (m *visitor) LogEvent(ctx context.Context, op scop.LogEvent) error {
	descID := screl.GetDescID(op.Element.Element())
	fullName, err := m.nr.GetFullyQualifiedName(ctx, descID)
	if err != nil {
		return err
	}
	event, err := asEventPayload(ctx, fullName, op.Element.Element(), op.TargetStatus, m)
	if err != nil {
		return err
	}
	details := eventpb.CommonSQLEventDetails{
		ApplicationName: op.Authorization.AppName,
		User:            op.Authorization.UserName,
		Statement:       redact.RedactableString(op.Statement),
		Tag:             op.StatementTag,
	}
	return m.s.EnqueueEvent(descID, op.TargetMetadata, details, event)
}

func asEventPayload(
	ctx context.Context, fullName string, e scpb.Element, targetStatus scpb.Status, m *visitor,
) (eventpb.EventPayload, error) {
	if targetStatus == scpb.Status_ABSENT {
		switch e.(type) {
		case *scpb.Table:
			return &eventpb.DropTable{TableName: fullName}, nil
		case *scpb.View:
			return &eventpb.DropView{ViewName: fullName}, nil
		case *scpb.Sequence:
			return &eventpb.DropSequence{SequenceName: fullName}, nil
		case *scpb.Database:
			return &eventpb.DropDatabase{DatabaseName: fullName}, nil
		case *scpb.Schema:
			return &eventpb.DropSchema{SchemaName: fullName}, nil
		case *scpb.AliasType, *scpb.EnumType:
			return &eventpb.DropType{TypeName: fullName}, nil
		}
	}
	switch e := e.(type) {
	case *scpb.Column:
		tbl, err := m.checkOutTable(ctx, e.TableID)
		if err != nil {
			return nil, err
		}
		mutation, err := FindMutation(tbl, MakeColumnIDMutationSelector(e.ColumnID))
		if err != nil {
			return nil, err
		}
		return &eventpb.AlterTable{
			TableName:  fullName,
			MutationID: uint32(mutation.MutationID()),
		}, nil
	case *scpb.SecondaryIndex:
		tbl, err := m.checkOutTable(ctx, e.TableID)
		if err != nil {
			return nil, err
		}
		mutation, err := FindMutation(tbl, MakeIndexIDMutationSelector(e.IndexID))
		if err != nil {
			return nil, err
		}
		switch targetStatus {
		case scpb.Status_PUBLIC:
			return &eventpb.AlterTable{
				TableName:  fullName,
				MutationID: uint32(mutation.MutationID()),
			}, nil
		case scpb.Status_ABSENT:
			return &eventpb.DropIndex{
				TableName:  fullName,
				IndexName:  mutation.AsIndex().GetName(),
				MutationID: uint32(mutation.MutationID()),
			}, nil
		default:
			return nil, errors.AssertionFailedf("unknown target status %s", targetStatus)
		}
	}
	return nil, errors.AssertionFailedf("unknown %s element type %T", targetStatus.String(), e)
}
