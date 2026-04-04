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

// Package migration captures the facilities needed to define and execute
// migrations for a crdb cluster. These migrations can be arbitrarily long
// running, are free to send out arbitrary requests cluster wide, change
// internal DB state, and much more. They're typically reserved for crdb
// internal operations and state. Each migration is idempotent in nature, is
// associated with a specific cluster version, and executed when the cluster
// version is made active on every node in the cluster.
//
// Examples of migrations that apply would be migrations to move all raft state
// from one storage engine to another, or purging all usage of the replicated
// truncated state in KV.
package migration

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
)

// Migration defines a program to be executed once every node in the cluster is
// (a) running a specific binary version, and (b) has completed all prior
// migrations. Note that there are two types of migrations, a SystemMigration
// and a TenantMigration. A SystemMigration only runs on the system tenant and
// is used to migrate state at the KV layer. A TenantMigration runs on all
// tenants (including the system tenant) and should be used whenever state at
// the SQL layer is being migrated.
//
// Each migration is associated with a specific internal cluster version and is
// idempotent in nature. When setting the cluster version (via `SET CLUSTER
// SETTING version`), the manager process determines the set of migrations
// needed to bridge the gap between the current active cluster version, and the
// target one. See [1] for where that happens.
//
// To introduce a migration, start by adding version key to pkg/clusterversion
// and introducing a corresponding internal cluster version for it. See [2] for
// more details. Following that, define a Migration in the migrations package
// and add it to the appropriate migrations slice to the registry. Be sure to
// key it in with the new cluster version we just added. During cluster
// upgrades, once the operator is able to set a cluster version setting that's
// past the version that was introduced (typically the major release version
// the migration was introduced in), the manager will execute the defined
// migration before letting the upgrade finalize.
//
// If the migration requires below-Raft level changes ([3] is one example),
// you'll need to add a version switch and the relevant system-level migration
// in [4]. See IterateRangeDescriptors and the Migrate KV request for more
// details.
//
// [1]: `(*Manager).Migrate`
// [2]: pkg/clusterversion/cockroach_versions.go
// [3]: truncatedStateMigration
// [4]: pkg/kv/kvserver/batch_eval/cmd_migrate.go
//
type Migration interface {
	ClusterVersion() clusterversion.ClusterVersion
	Name() string
	internal() // restrict implementations to this package
}

// JobDeps are migration-specific dependencies used by the migration job to run
// migrations.
type JobDeps interface {

	// GetMigration returns the migration associated with the cluster version
	// if one exists.
	GetMigration(key clusterversion.ClusterVersion) (Migration, bool)

	// SystemDeps returns a handle to migration dependencies on a system tenant.
	SystemDeps() SystemDeps
}

type migration struct {
	description string
	cv          clusterversion.ClusterVersion
}

// ClusterVersion makes SystemMigration a Migration.
func (m *migration) ClusterVersion() clusterversion.ClusterVersion {
	return m.cv
}

// Name returns a human readable name for this migration.
func (m *migration) Name() string {
	return fmt.Sprintf("Migration to %s: %q", m.cv.String(), m.description)
}

func (m *migration) internal() {}
