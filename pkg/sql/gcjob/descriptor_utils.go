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

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descs"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// updateDescriptorGCMutations removes the GCMutation for the given
// index ID. We no longer populate this field, but we still search it
// to remove existing entries.
func updateDescriptorGCMutations(
	ctx context.Context,
	execCfg *sql.ExecutorConfig,
	tableID descpb.ID,
	garbageCollectedIndexID descpb.IndexID,
) error {
	return sql.DescsTxn(ctx, execCfg, func(
		ctx context.Context, txn *kv.Txn, descsCol *descs.Collection,
	) error {
		tbl, err := descsCol.GetMutableTableVersionByID(ctx, tableID, txn)
		if err != nil {
			return err
		}
		found := false
		for i := 0; i < len(tbl.GCMutations); i++ {
			other := tbl.GCMutations[i]
			if other.IndexID == garbageCollectedIndexID {
				tbl.GCMutations = append(tbl.GCMutations[:i], tbl.GCMutations[i+1:]...)
				found = true
				break
			}
		}
		if found {
			log.Infof(ctx, "updating GCMutations for table %d after removing index %d",
				tableID, garbageCollectedIndexID)
			// Remove the mutation from the table descriptor.
			b := txn.NewBatch()
			if err := descsCol.WriteDescToBatch(ctx, false /* kvTrace */, tbl, b); err != nil {
				return err
			}
			return txn.Run(ctx, b)
		}
		return nil
	})
}

// deleteDatabaseZoneConfig removes the zone config for a given database ID.
func deleteDatabaseZoneConfig(
	ctx context.Context,
	db *kv.DB,
	codec keys.SQLCodec,
	settings *cluster.Settings,
	databaseID descpb.ID,
) error {
	if databaseID == descpb.InvalidID {
		return nil
	}
	return db.Txn(ctx, func(ctx context.Context, txn *kv.Txn) error {
		if !settings.Version.IsActive(
			ctx, clusterversion.DisableSystemConfigGossipTrigger,
		) {
			if err := txn.DeprecatedSetSystemConfigTrigger(codec.ForSystemTenant()); err != nil {
				return err
			}
		}
		b := &kv.Batch{}

		// Delete the zone config entry for the dropped database associated with the
		// job, if it exists.
		dbZoneKeyPrefix := config.MakeZoneKeyPrefix(codec, databaseID)
		b.DelRange(dbZoneKeyPrefix, dbZoneKeyPrefix.PrefixEnd(), false /* returnKeys */)
		return txn.Run(ctx, b)
	})
}
