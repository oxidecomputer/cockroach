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

package gcjob

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/jobs/jobspb"
	"github.com/cockroachdb/cockroach/pkg/sql"
)

// TestingGCTenant is a wrapper around the internal function that gc-s a tenant
// made public for testing.
func TestingGCTenant(
	ctx context.Context,
	execCfg *sql.ExecutorConfig,
	tenID uint64,
	progress *jobspb.SchemaChangeGCProgress,
) error {
	return gcTenant(ctx, execCfg, tenID, progress)
}
