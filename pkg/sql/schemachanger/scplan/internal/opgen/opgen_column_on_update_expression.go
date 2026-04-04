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

package opgen

import (
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
)

func init() {
	opRegistry.register((*scpb.ColumnOnUpdateExpression)(nil),
		toPublic(
			scpb.Status_ABSENT,
			to(scpb.Status_PUBLIC,
				minPhase(scop.PreCommitPhase),
				emit(func(this *scpb.ColumnOnUpdateExpression) scop.Op {
					return &scop.AddColumnOnUpdateExpression{
						OnUpdate: *protoutil.Clone(this).(*scpb.ColumnOnUpdateExpression),
					}
				}),
				emit(func(this *scpb.ColumnOnUpdateExpression) scop.Op {
					if len(this.UsesTypeIDs) == 0 {
						return nil
					}
					return &scop.UpdateTableBackReferencesInTypes{
						TypeIDs:               this.UsesTypeIDs,
						BackReferencedTableID: this.TableID,
					}
				}),
				emit(func(this *scpb.ColumnOnUpdateExpression) scop.Op {
					if len(this.UsesSequenceIDs) == 0 {
						return nil
					}
					return &scop.UpdateBackReferencesInSequences{
						SequenceIDs:            this.UsesSequenceIDs,
						BackReferencedTableID:  this.TableID,
						BackReferencedColumnID: this.ColumnID,
					}
				}),
			),
		),
		toAbsent(
			scpb.Status_PUBLIC,
			to(scpb.Status_ABSENT,
				minPhase(scop.PreCommitPhase),
				emit(func(this *scpb.ColumnOnUpdateExpression) scop.Op {
					return &scop.RemoveColumnOnUpdateExpression{
						TableID:  this.TableID,
						ColumnID: this.ColumnID,
					}
				}),
				emit(func(this *scpb.ColumnOnUpdateExpression) scop.Op {
					if len(this.UsesTypeIDs) == 0 {
						return nil
					}
					return &scop.UpdateTableBackReferencesInTypes{
						TypeIDs:               this.UsesTypeIDs,
						BackReferencedTableID: this.TableID,
					}
				}),
				emit(func(this *scpb.ColumnOnUpdateExpression) scop.Op {
					if len(this.UsesSequenceIDs) == 0 {
						return nil
					}
					return &scop.UpdateBackReferencesInSequences{
						SequenceIDs:            this.UsesSequenceIDs,
						BackReferencedTableID:  this.TableID,
						BackReferencedColumnID: this.ColumnID,
					}
				}),
			),
		),
	)
}
