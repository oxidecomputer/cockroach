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
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
)

func init() {
	opRegistry.register((*scpb.TemporaryIndex)(nil),
		toPublic(
			scpb.Status_ABSENT,
			to(scpb.Status_DELETE_ONLY,
				minPhase(scop.PreCommitPhase),
				emit(func(this *scpb.TemporaryIndex) scop.Op {
					return &scop.MakeAddedIndexDeleteOnly{
						Index:              *protoutil.Clone(&this.Index).(*scpb.Index),
						IsSecondaryIndex:   this.IsUsingSecondaryEncoding,
						IsDeletePreserving: true,
					}
				}),
			),
			to(scpb.Status_WRITE_ONLY,
				minPhase(scop.PostCommitPhase),
				emit(func(this *scpb.TemporaryIndex) scop.Op {
					return &scop.MakeAddedIndexDeleteAndWriteOnly{
						TableID: this.TableID,
						IndexID: this.IndexID,
					}
				}),
			),
			to(scpb.Status_PUBLIC),
		),
		toAbsent(
			scpb.Status_PUBLIC,
			equiv(scpb.Status_WRITE_ONLY),
			to(scpb.Status_DELETE_ONLY,
				revertible(false),
				emit(func(this *scpb.TemporaryIndex) scop.Op {
					return &scop.MakeDroppedIndexDeleteOnly{
						TableID: this.TableID,
						IndexID: this.IndexID,
					}
				}),
			),
			to(scpb.Status_ABSENT,
				emit(func(this *scpb.TemporaryIndex) scop.Op {
					return &scop.CreateGcJobForIndex{
						TableID: this.TableID,
						IndexID: this.IndexID,
					}
				}),
				emit(func(this *scpb.TemporaryIndex) scop.Op {
					return &scop.MakeIndexAbsent{
						TableID: this.TableID,
						IndexID: this.IndexID,
					}
				}),
			),
		),
	)
}
