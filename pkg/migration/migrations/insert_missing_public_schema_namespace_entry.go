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
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catalogkeys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descs"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// insertMissingPublicSchemaNamespaceEntry creates a system.namespace entries
// for public schemas that are missing a system.namespace entry.
// This arises from restore where we mistakenly did not create system.namespace
// entries for public schemas when restoring databases.
func insertMissingPublicSchemaNamespaceEntry(
	ctx context.Context, _ clusterversion.ClusterVersion, d migration.TenantDeps, _ *jobs.Job,
) error {
	// Get the ID of all databases where we're missing a public schema namespace
	// entry for.
	query := `
  SELECT id
    FROM system.namespace
   WHERE id
         NOT IN (
                SELECT ns_db.id
                  FROM system.namespace AS ns_db
                       INNER JOIN system.namespace
                            AS ns_sc ON (
                                        ns_db.id
                                        = ns_sc."parentID"
                                    )
                 WHERE ns_db."parentSchemaID" = 0
                   AND ns_db."parentID" = 0
                   AND ns_sc."parentSchemaID" = 0
                   AND ns_sc.name = 'public'
                   AND ns_sc.id = 29
            )
     AND "parentID" = 0
ORDER BY id ASC;
`
	rows, err := d.InternalExecutor.QueryIterator(
		ctx, "get_databases_without_public_schema_namespace_entry", nil /* txn */, query,
	)
	if err != nil {
		return err
	}
	var databaseIDs []descpb.ID
	for ok, err := rows.Next(ctx); ok; ok, err = rows.Next(ctx) {
		if err != nil {
			return err
		}
		id := descpb.ID(tree.MustBeDInt(rows.Cur()[0]))
		databaseIDs = append(databaseIDs, id)
	}

	return d.CollectionFactory.Txn(ctx, d.InternalExecutor, d.DB, func(
		ctx context.Context, txn *kv.Txn, descriptors *descs.Collection,
	) error {
		b := txn.NewBatch()
		for _, dbID := range databaseIDs {
			publicSchemaKey := catalogkeys.MakeSchemaNameKey(d.Codec, dbID, tree.PublicSchema)
			b.Put(publicSchemaKey, keys.PublicSchemaID)
		}
		return txn.Run(ctx, b)
	})
}
