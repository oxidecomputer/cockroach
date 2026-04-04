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

package migration

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/jobs"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descs"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/lease"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlutil"
	"github.com/cockroachdb/logtags"
)

// TenantDeps are the dependencies of migrations which perform actions at the
// SQL layer.
type TenantDeps struct {
	DB                *kv.DB
	Codec             keys.SQLCodec
	Settings          *cluster.Settings
	CollectionFactory *descs.CollectionFactory
	LeaseManager      *lease.Manager
	InternalExecutor  sqlutil.InternalExecutor

	SpanConfig struct { // deps for span config migrations; can be removed accordingly
		spanconfig.KVAccessor
		spanconfig.Splitter
		Default roachpb.SpanConfig
	}

	TestingKnobs *TestingKnobs
}

// TenantMigrationFunc is used to perform sql-level migrations. It may be run from
// any tenant.
type TenantMigrationFunc func(context.Context, clusterversion.ClusterVersion, TenantDeps, *jobs.Job) error

// PreconditionFunc is a function run without isolation before attempting an
// upgrade that includes this migration. It is used to verify that the
// required conditions for the migration to succeed are met. This can allow
// users to fix any problems before "crossing the rubicon" and no longer
// being able to upgrade.
type PreconditionFunc func(context.Context, clusterversion.ClusterVersion, TenantDeps) error

// TenantMigration is an implementation of Migration for tenant-level
// migrations. This is used for all migration which might affect the state of
// sql. It includes the system tenant.
type TenantMigration struct {
	migration
	fn           TenantMigrationFunc
	precondition PreconditionFunc
}

var _ Migration = (*TenantMigration)(nil)

// NewTenantMigration constructs a TenantMigration.
func NewTenantMigration(
	description string,
	cv clusterversion.ClusterVersion,
	precondition PreconditionFunc,
	fn TenantMigrationFunc,
) *TenantMigration {
	m := &TenantMigration{
		migration: migration{
			description: description,
			cv:          cv,
		},
		fn:           fn,
		precondition: precondition,
	}
	return m
}

// Run kickstarts the actual migration process for tenant-level migrations.
func (m *TenantMigration) Run(
	ctx context.Context, cv clusterversion.ClusterVersion, d TenantDeps, job *jobs.Job,
) error {
	ctx = logtags.AddTag(ctx, fmt.Sprintf("migration=%s", cv), nil)
	return m.fn(ctx, cv, d, job)
}

// Precondition runs the precondition check if there is one and reports
// any errors.
func (m *TenantMigration) Precondition(
	ctx context.Context, cv clusterversion.ClusterVersion, d TenantDeps,
) error {
	ctx = logtags.AddTag(ctx, fmt.Sprintf("migration=%s,precondition", cv), nil)
	return m.precondition(ctx, cv, d)
}
