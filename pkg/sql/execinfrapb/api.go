// Copyright 2018 The Cockroach Authors.
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

package execinfrapb

import (
	"strconv"

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
)

// ProcessorID identifies a processor in the context of a specific flow.
type ProcessorID int

// StreamID identifies a stream; it may be local to a flow or it may cross
// machine boundaries. The identifier can only be used in the context of a
// specific flow.
type StreamID int

func (sid StreamID) String() string {
	return strconv.Itoa(int(sid))
}

// FlowID identifies a flow. It is most importantly used when setting up streams
// between nodes.
type FlowID struct {
	uuid.UUID
}

// Equal returns whether the two FlowIDs are equal.
func (f FlowID) Equal(other FlowID) bool {
	return f.UUID.Equal(other.UUID)
}

// IsUnset returns whether the FlowID is unset.
func (f FlowID) IsUnset() bool {
	return f.UUID.Equal(uuid.Nil)
}

// DistSQLVersion identifies DistSQL engine versions.
type DistSQLVersion uint32

// MakeEvalContext serializes some of the fields of a tree.EvalContext into a
// execinfrapb.EvalContext proto.
func MakeEvalContext(evalCtx *tree.EvalContext) EvalContext {
	sessionDataProto := evalCtx.SessionData().SessionData
	sessiondata.MarshalNonLocal(evalCtx.SessionData(), &sessionDataProto)
	return EvalContext{
		SessionData:        sessionDataProto,
		StmtTimestampNanos: evalCtx.StmtTimestamp.UnixNano(),
		TxnTimestampNanos:  evalCtx.TxnTimestamp.UnixNano(),
	}
}

// User accesses the user field.
func (m *BackupDataSpec) User() security.SQLUsername {
	return m.UserProto.Decode()
}

// User accesses the user field.
func (m *ExportSpec) User() security.SQLUsername {
	return m.UserProto.Decode()
}

// User accesses the user field.
func (m *ReadImportDataSpec) User() security.SQLUsername {
	return m.UserProto.Decode()
}

// User accesses the user field.
func (m *ChangeAggregatorSpec) User() security.SQLUsername {
	return m.UserProto.Decode()
}

// User accesses the user field.
func (m *ChangeFrontierSpec) User() security.SQLUsername {
	return m.UserProto.Decode()
}
