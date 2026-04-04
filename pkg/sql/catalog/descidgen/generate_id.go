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

package descidgen

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
)

// GenerateUniqueDescID returns the next available Descriptor ID and increments
// the counter. The incrementing is non-transactional, and the counter could be
// incremented multiple times because of retries.
//
// TODO(postamar): define an interface for ID generation,
// both the db and codec dependencies are singletons anyway.
func GenerateUniqueDescID(ctx context.Context, db *kv.DB, codec keys.SQLCodec) (descpb.ID, error) {
	// Increment unique descriptor counter.
	newVal, err := kv.IncrementValRetryable(ctx, db, codec.DescIDSequenceKey(), 1)
	if err != nil {
		return descpb.InvalidID, err
	}
	return descpb.ID(newVal - 1), nil
}

// PeekNextUniqueDescID returns the next as-of-yet unassigned unique descriptor
// ID in the sequence. Note that this value is _not_ guaranteed to be the same
// as that returned by a subsequent call to GenerateUniqueDescID. It will,
// however, be a lower bound on it.
func PeekNextUniqueDescID(ctx context.Context, db *kv.DB, codec keys.SQLCodec) (descpb.ID, error) {
	v, err := db.Get(ctx, codec.DescIDSequenceKey())
	if err != nil {
		return descpb.InvalidID, err
	}
	return descpb.ID(v.ValueInt()), nil
}
