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

package systemschema

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/dbdesc"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/schemadesc"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/tabledesc"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/typedesc"
	"github.com/stretchr/testify/require"
)

func TestShouldSplitAtDesc(t *testing.T) {
	tbl1 := descpb.TableDescriptor{}
	tbl2 := descpb.TableDescriptor{ViewQuery: "SELECT"}
	tbl3 := descpb.TableDescriptor{ViewQuery: "SELECT", IsMaterializedView: true}
	typ := descpb.TypeDescriptor{}
	schema := descpb.SchemaDescriptor{}
	for inner, should := range map[catalog.Descriptor]bool{
		tabledesc.NewBuilder(&tbl1).BuildImmutable():          true,
		tabledesc.NewBuilder(&tbl2).BuildImmutable():          false,
		tabledesc.NewBuilder(&tbl3).BuildImmutable():          true,
		dbdesc.NewInitial(42, "db", security.AdminRoleName()): false,
		typedesc.NewBuilder(&typ).BuildCreatedMutable():       false,
		schemadesc.NewBuilder(&schema).BuildImmutable():       false,
	} {
		var rawDesc roachpb.Value
		require.NoError(t, rawDesc.SetProto(inner.DescriptorProto()))
		require.Equal(t, should, ShouldSplitAtDesc(&rawDesc))
	}
}
