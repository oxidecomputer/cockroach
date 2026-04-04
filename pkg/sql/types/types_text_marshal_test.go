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

package types_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/typedesc"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/stretchr/testify/require"
)

// TestDescriptorProtoString is to make sure gogo/protobuf is able to text
// marshal a protobuf struct has child field of type EnumMetadata
func TestDescriptorProtoString(t *testing.T) {
	enumMembers := []string{"hi", "hello"}
	enumType := types.MakeEnum(typedesc.TypeIDToOID(500), typedesc.TypeIDToOID(100500))
	enumType.TypeMeta = types.UserDefinedTypeMetadata{
		Name: &types.UserDefinedTypeName{
			Schema: "test",
			Name:   "greeting",
		},
		EnumData: &types.EnumMetadata{
			LogicalRepresentations: enumMembers,
			PhysicalRepresentations: [][]byte{
				{0x42, 0x1},
				{0x42},
			},
			IsMemberReadOnly: make([]bool, len(enumMembers)),
		},
	}
	desc := &descpb.ColumnDescriptor{
		Name: "c",
		ID:   1,
		Type: enumType,
	}

	var str string
	require.NotPanics(t, func() { str = desc.String() })
	// Assert we only dump InternalType from types.T without metadata
	require.Contains(t, str, "type:<family")
	require.NotContains(t, str, "TypeMeta:<Name")
}
