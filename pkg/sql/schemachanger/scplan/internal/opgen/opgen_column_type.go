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
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/catid"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
)

func init() {
	opRegistry.register((*scpb.ColumnType)(nil),
		toPublic(
			scpb.Status_ABSENT,
			to(scpb.Status_PUBLIC,
				minPhase(scop.PreCommitPhase),
				emit(func(this *scpb.ColumnType) scop.Op {
					return &scop.SetAddedColumnType{
						ColumnType: *protoutil.Clone(this).(*scpb.ColumnType),
					}
				}),
				emit(func(this *scpb.ColumnType) scop.Op {
					if ids := referencedTypeIDs(this); len(ids) > 0 {
						return &scop.UpdateTableBackReferencesInTypes{
							TypeIDs:               ids,
							BackReferencedTableID: this.TableID,
						}
					}
					return nil
				}),
			),
		),
		toAbsent(
			scpb.Status_PUBLIC,
			to(scpb.Status_ABSENT,
				minPhase(scop.PreCommitPhase),
				revertible(false),
				emit(func(this *scpb.ColumnType) scop.Op {
					return &scop.RemoveDroppedColumnType{
						TableID:  this.TableID,
						ColumnID: this.ColumnID,
					}
				}),
				emit(func(this *scpb.ColumnType) scop.Op {
					if ids := referencedTypeIDs(this); len(ids) > 0 {
						return &scop.UpdateTableBackReferencesInTypes{
							TypeIDs:               ids,
							BackReferencedTableID: this.TableID,
						}
					}
					return nil
				}),
			),
		),
	)
}

func referencedTypeIDs(this *scpb.ColumnType) []catid.DescID {
	var ids catalog.DescriptorIDSet
	if this.ComputeExpr != nil {
		for _, id := range this.ComputeExpr.UsesTypeIDs {
			ids.Add(id)
		}
	}
	for _, id := range this.ClosedTypeIDs {
		ids.Add(id)
	}
	return ids.Ordered()
}
