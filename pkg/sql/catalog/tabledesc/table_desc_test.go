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

package tabledesc

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/stretchr/testify/require"
)

func TestMaybeIncrementVersion(t *testing.T) {
	// Created descriptors should not have their version incremented.
	t.Run("created does not get incremented", func(t *testing.T) {
		{
			mut := NewBuilder(&descpb.TableDescriptor{
				ID:      1,
				Version: 1,
			}).BuildCreatedMutableTable()
			mut.MaybeIncrementVersion()
			require.Equal(t, descpb.DescriptorVersion(1), mut.GetVersion())
		}
		{
			mut := NewBuilder(&descpb.TableDescriptor{
				ID:      1,
				Version: 42,
			}).BuildCreatedMutableTable()
			mut.MaybeIncrementVersion()
			require.Equal(t, descpb.DescriptorVersion(42), mut.GetVersion())
		}
	})
	t.Run("existed gets incremented once", func(t *testing.T) {
		mut := NewBuilder(&descpb.TableDescriptor{
			ID:      1,
			Version: 1,
		}).BuildExistingMutableTable()
		require.Equal(t, descpb.DescriptorVersion(1), mut.GetVersion())
		mut.MaybeIncrementVersion()
		require.Equal(t, descpb.DescriptorVersion(2), mut.GetVersion())
		mut.MaybeIncrementVersion()
		require.Equal(t, descpb.DescriptorVersion(2), mut.GetVersion())
	})
}

// TestingSetClusterVersion is a test helper to override the original table
// descriptor.
func (desc *Mutable) TestingSetClusterVersion(d descpb.TableDescriptor) {
	desc.original = makeImmutable(&d)
}
