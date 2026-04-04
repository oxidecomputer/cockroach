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

package migrations

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/migration"
)

// runRemoveInvalidDatabasePrivileges calls RunPostDeserializationChanges on
// every database descriptor. It also calls RunPostDeserializationChanges on
// all table descriptors to add constraint IDs.
// This migration is done to convert invalid privileges on the
// database to default privileges.
func runRemoveInvalidDatabasePrivileges(
	ctx context.Context, _ clusterversion.ClusterVersion, d migration.TenantDeps, _ *jobs.Job,
) error {
	return runPostDeserializationChangesOnAllDescriptors(ctx, d)
}
