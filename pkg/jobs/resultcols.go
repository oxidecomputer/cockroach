// Copyright 2020 The Cockroach Authors.
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

package jobs

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// BulkJobExecutionResultHeader is the header for various job commands
// (BACKUP, RESTORE, IMPORT, etc) stmt results.
var BulkJobExecutionResultHeader = colinfo.ResultColumns{
	{Name: "job_id", Typ: types.Int},
	{Name: "status", Typ: types.String},
	{Name: "fraction_completed", Typ: types.Float},
	{Name: "rows", Typ: types.Int},
	{Name: "index_entries", Typ: types.Int},
	{Name: "bytes", Typ: types.Int},
}

// DetachedJobExecutionResultHeader is a the header for various job commands when
// job executes in detached mode (i.e. the caller doesn't wait for job completion).
var DetachedJobExecutionResultHeader = colinfo.ResultColumns{
	{Name: "job_id", Typ: types.Int},
}
