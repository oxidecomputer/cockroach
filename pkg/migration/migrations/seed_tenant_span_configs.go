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
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/migration"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/errors"
)

func seedTenantSpanConfigsMigration(
	ctx context.Context, _ clusterversion.ClusterVersion, d migration.TenantDeps, _ *jobs.Job,
) error {
	if !d.Codec.ForSystemTenant() {
		return nil
	}

	return d.DB.Txn(ctx, func(ctx context.Context, txn *kv.Txn) error {
		const getTenantIDsQuery = `SELECT id from system.tenants`
		it, err := d.InternalExecutor.QueryIteratorEx(ctx, "get-tenant-ids", txn,
			sessiondata.NodeUserSessionDataOverride, getTenantIDsQuery,
		)
		if err != nil {
			return errors.Wrap(err, "unable to fetch existing tenant IDs")
		}

		var tenantIDs []roachpb.TenantID
		var ok bool
		for ok, err = it.Next(ctx); ok; ok, err = it.Next(ctx) {
			row := it.Cur()
			tenantID := roachpb.MakeTenantID(uint64(tree.MustBeDInt(row[0])))
			tenantIDs = append(tenantIDs, tenantID)
		}
		if err != nil {
			return err
		}

		scKVAccessor := d.SpanConfig.KVAccessor.WithTxn(ctx, txn)
		for _, tenantID := range tenantIDs {
			// Install a single key span config at the start of tenant's
			// keyspace; elsewhere this ensures that we split on the tenant
			// boundary. Look towards CreateTenantRecord for more details.
			tenantSpanConfig := d.SpanConfig.Default
			tenantPrefix := keys.MakeTenantPrefix(tenantID)
			tenantTarget := spanconfig.MakeTargetFromSpan(roachpb.Span{
				Key:    tenantPrefix,
				EndKey: tenantPrefix.PrefixEnd(),
			})
			tenantSeedSpan := roachpb.Span{
				Key:    tenantPrefix,
				EndKey: tenantPrefix.Next(),
			}
			record, err := spanconfig.MakeRecord(spanconfig.MakeTargetFromSpan(tenantSeedSpan),
				tenantSpanConfig)
			if err != nil {
				return err
			}
			toUpsert := []spanconfig.Record{record}
			scRecords, err := scKVAccessor.GetSpanConfigRecords(ctx, []spanconfig.Target{tenantTarget})
			if err != nil {
				return err
			}
			if len(scRecords) != 0 {
				// This tenant already has span config records. It was either
				// already migrated (migrations need to be idempotent) or it was
				// created after PreSeedTenantSpanConfigs was activated. There's
				// nothing left to do here.
				continue
			}
			if err := scKVAccessor.UpdateSpanConfigRecords(
				ctx, nil /* toDelete */, toUpsert, hlc.MinTimestamp, hlc.MaxTimestamp,
			); err != nil {
				return errors.Wrapf(err, "failed to seed span config for tenant %d", tenantID)
			}
		}

		return nil
	})
}
