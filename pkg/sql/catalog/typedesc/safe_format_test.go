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

package typedesc_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/typedesc"
	"github.com/cockroachdb/redact"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSafeMessage(t *testing.T) {
	for _, tc := range []struct {
		desc catalog.TypeDescriptor
		exp  string
	}{
		{
			desc: typedesc.NewBuilder(&descpb.TypeDescriptor{
				Name:                     "foo",
				ID:                       21,
				Version:                  3,
				Privileges:               catpb.NewBasePrivilegeDescriptor(security.RootUserName()),
				ParentID:                 2,
				ParentSchemaID:           29,
				ArrayTypeID:              117,
				State:                    descpb.DescriptorState_PUBLIC,
				Kind:                     descpb.TypeDescriptor_ALIAS,
				ReferencingDescriptorIDs: []descpb.ID{73, 37},
			}).BuildImmutableType(),
			exp: `typedesc.immutable: {ID: 21, Version: 3, ModificationTime: "0,0", ` +
				`ParentID: 2, ParentSchemaID: 29, State: PUBLIC, ` +
				`Kind: ALIAS, ArrayTypeID: 117, ReferencingDescriptorIDs: [73, 37]}`,
		},
		{
			desc: typedesc.NewBuilder(&descpb.TypeDescriptor{
				Name:                     "foo",
				ID:                       21,
				Version:                  3,
				Privileges:               catpb.NewBasePrivilegeDescriptor(security.RootUserName()),
				ParentID:                 2,
				ParentSchemaID:           29,
				ArrayTypeID:              117,
				State:                    descpb.DescriptorState_PUBLIC,
				Kind:                     descpb.TypeDescriptor_ENUM,
				ReferencingDescriptorIDs: []descpb.ID{73, 37},
				EnumMembers: []descpb.TypeDescriptor_EnumMember{
					{},
				},
			}).BuildImmutableType(),
			exp: `typedesc.immutable: {ID: 21, Version: 3, ModificationTime: "0,0", ` +
				`ParentID: 2, ParentSchemaID: 29, State: PUBLIC, ` +
				`Kind: ENUM, NumEnumMembers: 1, ArrayTypeID: 117, ReferencingDescriptorIDs: [73, 37]}`,
		},
	} {
		t.Run("", func(t *testing.T) {
			redacted := string(redact.Sprint(tc.desc).Redact())
			require.Equal(t, tc.exp, redacted)
			{
				var m map[string]interface{}
				require.NoError(t, yaml.UnmarshalStrict([]byte(redacted), &m))
			}
		})
	}
}
