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

package sql

import (
	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondatapb"
)

// SchemaChangerState is state associated with the new schema changer.
type SchemaChangerState struct {
	mode  sessiondatapb.NewSchemaChangerMode
	state scpb.CurrentState
	// jobID contains the ID of the schema changer job, if it is to be created.
	jobID jobspb.JobID
	// stmts contains the SQL statements involved in the schema change. This is
	// the bare minimum of statement information we need for testing, but in the
	// future we may want sql.Statement or something.
	stmts []string
}
