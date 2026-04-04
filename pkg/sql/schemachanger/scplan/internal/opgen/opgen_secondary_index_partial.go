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
)

func init() {
	opRegistry.register((*scpb.SecondaryIndexPartial)(nil),
		toPublic(
			scpb.Status_ABSENT,
			to(scpb.Status_PUBLIC,
				minPhase(scop.PreCommitPhase),
				emit(func(this *scpb.SecondaryIndexPartial) scop.Op {
					return &scop.SetAddedIndexPartialPredicate{
						TableID: this.TableID,
						IndexID: this.IndexID,
						Expr:    this.Expr,
					}
				}),
				emit(func(this *scpb.SecondaryIndexPartial) scop.Op {
					if len(this.UsesTypeIDs) == 0 {
						return nil
					}
					return &scop.UpdateTableBackReferencesInTypes{
						TypeIDs:               this.UsesTypeIDs,
						BackReferencedTableID: this.TableID,
					}
				}),
			),
		),
		toAbsent(
			scpb.Status_PUBLIC,
			to(scpb.Status_ABSENT,
				minPhase(scop.PreCommitPhase),
				revertible(false),
				emit(func(this *scpb.SecondaryIndexPartial) scop.Op {
					return &scop.RemoveDroppedIndexPartialPredicate{
						TableID: this.TableID,
						IndexID: this.IndexID,
					}
				}),
				emit(func(this *scpb.SecondaryIndexPartial) scop.Op {
					if len(this.UsesTypeIDs) == 0 {
						return nil
					}
					return &scop.UpdateTableBackReferencesInTypes{
						TypeIDs:               this.UsesTypeIDs,
						BackReferencedTableID: this.TableID,
					}
				}),
			),
		),
	)
}
