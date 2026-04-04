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

package opgen

import (
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

func init() {
	opRegistry.register((*scpb.CheckConstraint)(nil),
		toPublic(
			scpb.Status_ABSENT,
			to(scpb.Status_PUBLIC,
				emit(func(this *scpb.CheckConstraint) scop.Op {
					return notImplemented(this)
				}),
			),
		),
		toAbsent(
			scpb.Status_PUBLIC,
			to(scpb.Status_ABSENT,
				minPhase(scop.PreCommitPhase),
				// TODO(postamar): remove revertibility constraint when possible
				revertible(false),
				emit(func(this *scpb.CheckConstraint) scop.Op {
					return &scop.RemoveCheckConstraint{
						TableID:      this.TableID,
						ConstraintID: this.ConstraintID,
					}
				}),
				emit(func(this *scpb.CheckConstraint) scop.Op {
					if len(this.UsesTypeIDs) == 0 {
						return nil
					}
					return &scop.UpdateTableBackReferencesInTypes{
						TypeIDs:               this.UsesTypeIDs,
						BackReferencedTableID: this.TableID,
					}
				}),
				emit(func(this *scpb.CheckConstraint) scop.Op {
					if len(this.UsesSequenceIDs) == 0 {
						return nil
					}
					return &scop.UpdateBackReferencesInSequences{
						SequenceIDs:           this.UsesSequenceIDs,
						BackReferencedTableID: this.TableID,
					}
				}),
			),
		),
	)
}
