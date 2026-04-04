// Copyright 2019 The Cockroach Authors.
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

package sql

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestDefaultOid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	testCases := []struct {
		id  descpb.ID
		oid *tree.DOid
	}{
		{
			1,
			tree.NewDOid(tree.DInt(1)),
		},
		{
			2,
			tree.NewDOid(tree.DInt(2)),
		},
	}

	for _, tc := range testCases {
		oid := tableOid(tc.id)
		if tc.oid.DInt != oid.DInt {
			t.Fatalf("expected oid %d(%32b), got %d(%32b)", tc.oid.DInt, tc.oid.DInt, oid.DInt, oid.DInt)
		}
	}
}
