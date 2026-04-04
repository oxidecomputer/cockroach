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

package catkv

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catalogkeys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/nstree"
)

// LookupIDs returns the IDs of the descriptor for the requested namespace table
// row keys.
// descpb.InvalidID is returned for each request for which there exists no
// matching row in the namespace table.
func LookupIDs(
	ctx context.Context, txn *kv.Txn, codec keys.SQLCodec, nameInfos []descpb.NameInfo,
) ([]descpb.ID, error) {
	return lookupIDs(ctx, txn, catalogQuerier{codec: codec}, nameInfos)
}

// LookupID is like LookupIDs but for one record.
func LookupID(
	ctx context.Context,
	txn *kv.Txn,
	codec keys.SQLCodec,
	parentID descpb.ID,
	parentSchemaID descpb.ID,
	name string,
) (descpb.ID, error) {
	nameInfo := descpb.NameInfo{ParentID: parentID, ParentSchemaID: parentSchemaID, Name: name}
	ids, err := LookupIDs(ctx, txn, codec, []descpb.NameInfo{nameInfo})
	if err != nil {
		return descpb.InvalidID, err
	}
	return ids[0], nil
}

// GetAllDatabaseDescriptorIDs looks up and returns all available database
// descriptor IDs.
func GetAllDatabaseDescriptorIDs(
	ctx context.Context, txn *kv.Txn, codec keys.SQLCodec,
) (nstree.Catalog, error) {
	cq := catalogQuerier{codec: codec}
	return cq.query(ctx, txn, func(codec keys.SQLCodec, b *kv.Batch) {
		b.Header.MaxSpanRequestKeys = 0
		prefix := catalogkeys.MakeDatabaseNameKey(codec, "")
		b.Scan(prefix, prefix.PrefixEnd())
	})
}
