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

package opttester_test

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/testutils/opttester"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/testutils/testcat"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/datadriven"
)

func TestOptTester(t *testing.T) {
	const fmtFlags = memo.ExprFmtHideStats | memo.ExprFmtHideCost | memo.ExprFmtHideRuleProps |
		memo.ExprFmtHideQualifications

	datadriven.Walk(t, testutils.TestDataPath(t), func(t *testing.T, path string) {
		if strings.HasSuffix(path, ".json") {
			// Skip any .json files; used for inject-stats
			return
		}
		catalog := testcat.New()
		datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
			tester := opttester.New(catalog, d.Input)
			tester.Flags.ExprFormat = fmtFlags
			return tester.RunCommand(t, d)
		})
	})
}
