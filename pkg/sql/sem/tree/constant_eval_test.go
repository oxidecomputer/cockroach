// Copyright 2018 The Cockroach Authors.
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

package tree_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree/treecmp"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestConstantEvalArrayComparison(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	defer tree.MockNameTypes(map[string]*types.T{"a": types.MakeArray(types.Int)})()

	expr, err := parser.ParseExpr("a = ARRAY[1:::INT,2:::INT]")
	if err != nil {
		t.Fatal(err)
	}

	semaCtx := tree.MakeSemaContext()
	typedExpr, err := expr.TypeCheck(context.Background(), &semaCtx, types.Any)
	if err != nil {
		t.Fatal(err)
	}

	ctx := tree.NewTestingEvalContext(cluster.MakeTestingClusterSettings())
	defer ctx.Mon.Stop(context.Background())
	c := tree.MakeConstantEvalVisitor(ctx)
	expr, _ = tree.WalkExpr(&c, typedExpr)
	if err := c.Err(); err != nil {
		t.Fatal(err)
	}

	left := tree.ColumnItem{
		ColumnName: "a",
	}
	right := tree.DArray{
		ParamTyp:    types.Int,
		Array:       tree.Datums{tree.NewDInt(1), tree.NewDInt(2)},
		HasNonNulls: true,
	}
	expected := tree.NewTypedComparisonExpr(treecmp.MakeComparisonOperator(treecmp.EQ), &left, &right)
	if !reflect.DeepEqual(expr, expected) {
		t.Errorf("invalid expr '%v', expected '%v'", expr, expected)
	}
}
