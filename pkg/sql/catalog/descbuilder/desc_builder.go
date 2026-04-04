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

package descbuilder

import (
	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/dbdesc"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/internal/validate"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/schemadesc"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/tabledesc"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/typedesc"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
)

// NewBuilderWithMVCCTimestamp takes a descriptor as deserialized from storage,
// along with its MVCC timestamp, and returns a catalog.DescriptorBuilder object.
// Returns nil if nothing specific is found in desc.
func NewBuilderWithMVCCTimestamp(
	desc *descpb.Descriptor, mvccTimestamp hlc.Timestamp,
) catalog.DescriptorBuilder {
	table, database, typ, schema := descpb.FromDescriptorWithMVCCTimestamp(desc, mvccTimestamp)
	switch {
	case table != nil:
		return tabledesc.NewBuilder(table)
	case database != nil:
		return dbdesc.NewBuilder(database)
	case typ != nil:
		return typedesc.NewBuilder(typ)
	case schema != nil:
		return schemadesc.NewBuilder(schema)
	default:
		return nil
	}
}

// NewBuilder is a convenience function which calls NewBuilderWithMVCCTimestamp
// with an empty timestamp.
// The expectation here, therefore, is that this function is called either for a
// new descriptor that doesn't exist in storage yet, or already has its
// modification time field (and others which depend on the MVCC timestamp)
// set to a valid value.
func NewBuilder(desc *descpb.Descriptor) catalog.DescriptorBuilder {
	return NewBuilderWithMVCCTimestamp(desc, hlc.Timestamp{})
}

// ValidateSelf validates that the descriptor is internally consistent.
func ValidateSelf(desc catalog.Descriptor, version clusterversion.ClusterVersion) error {
	return validate.Self(version, desc)
}
