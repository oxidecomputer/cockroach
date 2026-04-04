// Copyright 2022 The Cockroach Authors.
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

package descs

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

// GetTableNameByID fetches the full tree table name by the given ID.
func GetTableNameByID(
	ctx context.Context, txn *kv.Txn, tc *Collection, tableID descpb.ID,
) (*tree.TableName, error) {
	tbl, err := tc.GetImmutableTableByID(ctx, txn, tableID, tree.ObjectLookupFlagsWithRequired())
	if err != nil {
		return nil, err
	}
	return GetTableNameByDesc(ctx, txn, tc, tbl)
}

// GetTableNameByDesc fetches the full tree table name by the given table descriptor.
func GetTableNameByDesc(
	ctx context.Context, txn *kv.Txn, tc *Collection, tbl catalog.TableDescriptor,
) (*tree.TableName, error) {
	sc, err := tc.GetImmutableSchemaByID(ctx, txn, tbl.GetParentSchemaID(), tree.SchemaLookupFlags{Required: true})
	if err != nil {
		return nil, err
	}
	found, db, err := tc.GetImmutableDatabaseByID(ctx, txn, tbl.GetParentID(), tree.DatabaseLookupFlags{Required: true})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.AssertionFailedf("expected database %d to exist", tbl.GetParentID())
	}
	return tree.NewTableNameWithSchema(tree.Name(db.GetName()), tree.Name(sc.GetName()), tree.Name(tbl.GetName())), nil
}
