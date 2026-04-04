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

	opRegistry.register((*scpb.Sequence)(nil),
		toPublic(
			scpb.Status_ABSENT,
			equiv(scpb.Status_TXN_DROPPED),
			equiv(scpb.Status_DROPPED),
			to(scpb.Status_PUBLIC,
				emit(func(this *scpb.Sequence) scop.Op {
					return notImplemented(this)
				}),
			),
		),
		toAbsent(scpb.Status_PUBLIC,
			to(scpb.Status_TXN_DROPPED,
				emit(func(this *scpb.Sequence) scop.Op {
					return &scop.MarkDescriptorAsDroppedSynthetically{
						DescID: this.SequenceID,
					}
				}),
			),
			to(scpb.Status_DROPPED,
				minPhase(scop.PreCommitPhase),
				revertible(false),
				emit(func(this *scpb.Sequence) scop.Op {
					return &scop.MarkDescriptorAsDropped{
						DescID: this.SequenceID,
					}
				}),
				emit(func(this *scpb.Sequence) scop.Op {
					return &scop.RemoveAllTableComments{
						TableID: this.SequenceID,
					}
				}),
			),
			to(scpb.Status_ABSENT,
				minPhase(scop.PostCommitPhase),
				emit(func(this *scpb.Sequence, md targetsWithElementMap) scop.Op {
					return newLogEventOp(this, md)
				}),
				emit(func(this *scpb.Sequence, md targetsWithElementMap) scop.Op {
					return &scop.CreateGcJobForTable{
						TableID:             this.SequenceID,
						StatementForDropJob: statementForDropJob(this, md),
					}
				}),
			),
		),
	)
}
