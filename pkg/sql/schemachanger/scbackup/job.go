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

package scbackup

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/nstree"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scexec"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/screl"
)

// CreateDeclarativeSchemaChangeJobs is called during the last phase of a
// restore. The provided catalog should contain all descriptors being restored.
// The code here will iterate those descriptors and synthesize the appropriate
// jobs.
//
// It should only be called for backups which do not restore the jobs table
// directly.
func CreateDeclarativeSchemaChangeJobs(
	ctx context.Context, registry *jobs.Registry, txn *kv.Txn, allMut nstree.Catalog,
) error {
	byJobID := make(map[catpb.JobID][]catalog.MutableDescriptor)
	_ = allMut.ForEachDescriptorEntry(func(d catalog.Descriptor) error {
		if s := d.GetDeclarativeSchemaChangerState(); s != nil {
			byJobID[s.JobID] = append(byJobID[s.JobID], d.(catalog.MutableDescriptor))
		}
		return nil
	})
	var records []*jobs.Record
	for _, descs := range byJobID {
		// TODO(ajwerner): Consider the need to trim elements or update
		// descriptors in the face of restoring only some constituent
		// descriptors of a larger change. One example where this needs
		// to happen urgently is sequences. Others shouldn't be possible
		// at this point.
		newID := registry.MakeJobID()
		var descriptorStates []*scpb.DescriptorState
		for _, d := range descs {
			ds := d.GetDeclarativeSchemaChangerState()
			ds.JobID = newID
			descriptorStates = append(descriptorStates, ds)
		}
		currentState, err := scpb.MakeCurrentStateFromDescriptors(
			descriptorStates,
		)
		if err != nil {
			return err
		}
		const runningStatus = "restored from backup"
		records = append(records, scexec.MakeDeclarativeSchemaChangeJobRecord(
			newID,
			currentState.Statements,
			!currentState.Revertible, // NonCancelable
			currentState.Authorization,
			screl.AllTargetDescIDs(currentState.TargetState).Ordered(),
			runningStatus,
		))
	}
	_, err := registry.CreateJobsWithTxn(ctx, txn, records)
	return err
}
