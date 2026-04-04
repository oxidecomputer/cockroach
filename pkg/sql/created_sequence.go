// Copyright 2016 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/errors"
)

type createdSequences interface {
	// addCreatedSequence adds a sequence to the set of sequences created in the current transaction.
	addCreatedSequence(id descpb.ID) error
	// isCreatedSequence checks if a sequence was created in the current transaction.
	isCreatedSequence(id descpb.ID) bool
}

type connExCreatedSequencesAccessor struct {
	ex *connExecutor
}

func (c connExCreatedSequencesAccessor) addCreatedSequence(id descpb.ID) error {
	c.ex.extraTxnState.createdSequences[id] = struct{}{}
	return nil
}

func (c connExCreatedSequencesAccessor) isCreatedSequence(id descpb.ID) bool {
	_, ok := c.ex.extraTxnState.createdSequences[id]
	return ok
}

// emptyCreatedSequences is the default impl used by the planner when the connExecutor is not available.
type emptyCreatedSequences struct{}

func (createdSequences emptyCreatedSequences) addCreatedSequence(id descpb.ID) error {
	return errors.AssertionFailedf("addCreatedSequence not supported in emptyCreatedSequences")
}

func (createdSequences emptyCreatedSequences) isCreatedSequence(id descpb.ID) bool {
	return false
}
