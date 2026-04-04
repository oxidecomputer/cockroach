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
	_ "github.com/cockroachdb/cockroach/pkg/sql/sem/builtins"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

// TestMarkRedactionStatement verifies that the redactable parts
// of statements are marked correctly.
func TestMarkRedactionStatement(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	testCases := []struct {
		query    string
		expected string
	}{
		{
			"SELECT 1",
			"SELECT ‹1›",
		},
		{
			"SELECT * FROM t",
			"SELECT * FROM ‹\"\"›.‹\"\"›.‹t›",
		},
		{
			"CREATE TABLE defaultdb.public.t()",
			"CREATE TABLE ‹defaultdb›.public.‹t› ()",
		},
		// Note: this test only passes due to the overriding of the FmtMarkRedactionNode flag
		// when formatting FuncExpr.Func. Although the overriding of the redaction flag
		// is correct (function names do not hold PII and therefore do not need redaction),
		// the incorrect typing of FuncExpr.Func as an UnresolvedName results in the
		// function name being incorrectly redacted when logging statements. In this
		// test case, that would be: SELECT ‹lower›(‹'foo'›).
		// The intended functionality is for function names to be correctly typed as
		// FunctionDefinition, which is logically identified as non-PII.
		{
			"SELECT lower('foo')",
			"SELECT lower(‹'foo'›)",
		},
		{
			"SELECT crdb_internal.node_executable_version()",
			"SELECT crdb_internal.node_executable_version()",
		},
		{
			"SHOW database",
			"SHOW database",
		},
		{
			"SET database = defaultdb",
			"SET database = ‹defaultdb›",
		},
		{
			"SET TRACING = off",
			"SET TRACING = off",
		},
		{
			"SHOW CLUSTER SETTING sql.log.admin_audit.enabled",
			"SHOW CLUSTER SETTING \"sql.log.admin_audit.enabled\"",
		},
		{
			"SET CLUSTER SETTING sql.log.admin_audit.enabled = true",
			"SET CLUSTER SETTING \"sql.log.admin_audit.enabled\" = true",
		},
	}

	s := cluster.MakeTestingClusterSettings()
	vt, err := NewVirtualSchemaHolder(context.Background(), s)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range testCases {
		stmt, err := parser.ParseOne(test.query)
		if err != nil {
			t.Fatal(err)
		}
		ann := tree.MakeAnnotations(stmt.NumAnnotations)
		f := tree.NewFmtCtx(
			tree.FmtAlwaysQualifyTableNames|tree.FmtMarkRedactionNode,
			tree.FmtAnnotations(&ann),
			tree.FmtReformatTableNames(hideNonVirtualTableNameFunc(vt)))
		f.FormatNode(stmt.AST)
		redactedString := f.CloseAndGetString()
		require.Equal(t, test.expected, redactedString)
	}
}
