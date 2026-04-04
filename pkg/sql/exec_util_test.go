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

package sql

import (
	"context"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestHideNonVirtualTableNameFunc(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	s := cluster.MakeTestingClusterSettings()
	vt, err := NewVirtualSchemaHolder(context.Background(), s)
	if err != nil {
		t.Fatal(err)
	}
	tableNameFunc := hideNonVirtualTableNameFunc(vt)

	testData := []struct {
		stmt     string
		expected string
	}{
		{`SELECT * FROM a.b`,
			`SELECT * FROM _`},
		{`SELECT * FROM pg_type`,
			`SELECT * FROM _`},
		{`SELECT * FROM pg_catalog.pg_type`,
			`SELECT * FROM pg_catalog.pg_type`},
	}

	for _, test := range testData {
		stmt, err := parser.ParseOne(test.stmt)
		if err != nil {
			t.Fatal(err)
		}
		f := tree.NewFmtCtx(
			tree.FmtSimple,
			tree.FmtReformatTableNames(tableNameFunc),
		)
		f.FormatNode(stmt.AST)
		actual := f.CloseAndGetString()
		require.Equal(t, test.expected, actual)
	}
}
