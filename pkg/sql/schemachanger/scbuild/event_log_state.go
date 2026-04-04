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

package scbuild

import (
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scbuild/internal/scbuildstmt"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

var _ scbuildstmt.EventLogState = (*eventLogState)(nil)

// TargetMetadata implements the scbuildstmt.EventLogState interface.
func (e *eventLogState) TargetMetadata() scpb.TargetMetadata {
	return e.statementMetaData
}

// IncrementSubWorkID implements the scbuildstmt.EventLogState interface.
func (e *eventLogState) IncrementSubWorkID() {
	e.statementMetaData.SubWorkID++
}

// EventLogStateWithNewSourceElementID implements the scbuildstmt.EventLogState
// interface.
func (e *eventLogState) EventLogStateWithNewSourceElementID() scbuildstmt.EventLogState {
	*e.sourceElementID++
	return &eventLogState{
		statements:      e.statements,
		authorization:   e.authorization,
		sourceElementID: e.sourceElementID,
		statementMetaData: scpb.TargetMetadata{
			StatementID:     e.statementMetaData.StatementID,
			SubWorkID:       e.statementMetaData.SubWorkID,
			SourceElementID: *e.sourceElementID,
		},
	}
}
