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

package norm_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/testutils/opttester"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/testutils/testcat"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/datadriven"
)

// TestNormRules tests the various Optgen normalization rules found in the rules
// directory. The tests are data-driven cases of the form:
//   <command>
//   <SQL statement>
//   ----
//   <expected results>
//
// See OptTester.Handle for supported commands.
//
// Rules files can be run separately like this:
//   make test PKG=./pkg/sql/opt/norm TESTS="TestNormRules/bool"
//   make test PKG=./pkg/sql/opt/norm TESTS="TestNormRules/comp"
//   ...
func TestNormRules(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	const fmtFlags = memo.ExprFmtHideStats | memo.ExprFmtHideCost | memo.ExprFmtHideRuleProps |
		memo.ExprFmtHideQualifications | memo.ExprFmtHideScalars | memo.ExprFmtHideTypes
	datadriven.Walk(t, testutils.TestDataPath(t, "rules"), func(t *testing.T, path string) {
		catalog := testcat.New()
		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			tester := opttester.New(catalog, d.Input)
			tester.Flags.ExprFormat = fmtFlags
			return tester.RunCommand(t, d)
		})
	})
}

// TestRuleProps files can be run separately like this:
//   make test PKG=./pkg/sql/opt/norm TESTS="TestNormRuleProps/orderings"
//   ...
func TestNormRuleProps(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	const fmtFlags = memo.ExprFmtHideStats | memo.ExprFmtHideCost |
		memo.ExprFmtHideQualifications | memo.ExprFmtHideScalars | memo.ExprFmtHideTypes
	datadriven.Walk(t, testutils.TestDataPath(t, "ruleprops"), func(t *testing.T, path string) {
		catalog := testcat.New()
		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			tester := opttester.New(catalog, d.Input)
			tester.Flags.ExprFormat = fmtFlags
			return tester.RunCommand(t, d)
		})
	})
}

// Ensure that every binary commutative operator overload can have its operands
// switched. Patterns like CommuteConst rely on this being possible.
func TestRuleBinaryAssumption(t *testing.T) {
	fn := func(op opt.Operator) {
		for _, overload := range tree.BinOps[opt.BinaryOpReverseMap[op]] {
			binOp := overload.(*tree.BinOp)
			if !memo.BinaryOverloadExists(op, binOp.RightType, binOp.LeftType) {
				t.Errorf("could not find inverse for overload: %+v", op)
			}
		}
	}

	// Only include commutative binary operators.
	fn(opt.PlusOp)
	fn(opt.MultOp)
	fn(opt.BitandOp)
	fn(opt.BitorOp)
	fn(opt.BitxorOp)
}
