// Copyright 2022 The Cockroach Authors.
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

package span_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/desctestutils"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/systemschema"
	"github.com/cockroachdb/cockroach/pkg/sql/span"
	"github.com/cockroachdb/cockroach/pkg/sql/tests"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestSpanSplitterDoesNotSplitSystemTableFamilySpans(t *testing.T) {
	splitter := span.MakeSplitter(
		systemschema.DescriptorTable,
		systemschema.DescriptorTable.GetPrimaryIndex(),
		util.MakeFastIntSet(0),
	)

	if res := splitter.CanSplitSpanIntoFamilySpans(1, false); res {
		t.Errorf("expected the system table to not be splittable")
	}
}

func TestSpanSplitterCanSplitSpan(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	ctx := context.Background()
	params, _ := tests.CreateTestServerParams()
	s, sqlDB, kvDB := serverutils.StartServer(t, params)
	defer s.Stopper().Stop(ctx)
	tcs := []struct {
		sql           string
		index         string
		prefixLen     int
		neededColumns util.FastIntSet
		containsNull  bool
		canSplit      bool
	}{
		{
			sql:           "a INT, b INT, c INT, d INT, PRIMARY KEY (a, b), FAMILY (a, b, c), FAMILY (d)",
			index:         "t_pkey",
			prefixLen:     2,
			neededColumns: util.MakeFastIntSet(0),
			canSplit:      true,
		},
		{
			sql:           "a INT, b INT, c INT, d INT, PRIMARY KEY (a, b), FAMILY (a, b, c), FAMILY (d)",
			index:         "t_pkey",
			prefixLen:     1,
			neededColumns: util.MakeFastIntSet(0),
			canSplit:      false,
		},
		{
			sql:           "a INT, b INT, c INT, d INT, PRIMARY KEY (a, b), FAMILY (a, b, c, d)",
			index:         "t_pkey",
			prefixLen:     2,
			neededColumns: util.MakeFastIntSet(0),
			canSplit:      true,
		},
		{
			sql:           "a INT, b INT, c INT, INDEX i (b) STORING (a, c), FAMILY (a), FAMILY (b), FAMILY (c)",
			index:         "i",
			prefixLen:     1,
			neededColumns: util.MakeFastIntSet(0),
			canSplit:      false,
		},
		{
			sql:           "a INT, b INT, c INT, UNIQUE INDEX i (b) STORING (a, c), FAMILY (a), FAMILY (b), FAMILY (c)",
			index:         "i",
			prefixLen:     1,
			neededColumns: util.MakeFastIntSet(0),
			containsNull:  true,
			canSplit:      false,
		},
		{
			sql:           "a INT, b INT, c INT, UNIQUE INDEX i (b) STORING (a, c), FAMILY (a), FAMILY (b), FAMILY (c)",
			index:         "i",
			prefixLen:     1,
			neededColumns: util.MakeFastIntSet(0),
			containsNull:  false,
			canSplit:      true,
		},
		{
			sql:           "a INT, b INT, c INT, UNIQUE INDEX i (b), FAMILY (a), FAMILY (b), FAMILY (c)",
			index:         "i",
			prefixLen:     1,
			neededColumns: util.MakeFastIntSet(0),
			containsNull:  false,
			canSplit:      true,
		},
		{
			sql:           "a INT, b INT, c INT, UNIQUE INDEX i (b), FAMILY (a), FAMILY (b), FAMILY (c)",
			index:         "i",
			prefixLen:     1,
			neededColumns: util.MakeFastIntSet(0),
			containsNull:  true,
			canSplit:      false,
		},
	}
	if _, err := sqlDB.Exec("CREATE DATABASE t"); err != nil {
		t.Fatal(err)
	}
	for _, tc := range tcs {
		t.Run(tc.sql, func(t *testing.T) {
			if _, err := sqlDB.Exec("DROP TABLE IF EXISTS t.t"); err != nil {
				t.Fatal(err)
			}
			sql := fmt.Sprintf("CREATE TABLE t.t (%s)", tc.sql)
			if _, err := sqlDB.Exec(sql); err != nil {
				t.Fatal(err)
			}
			desc := desctestutils.TestingGetPublicTableDescriptor(kvDB, keys.SystemSQLCodec, "t", "t")
			idx, err := desc.FindIndexWithName(tc.index)
			if err != nil {
				t.Fatal(err)
			}
			splitter := span.MakeSplitter(desc, idx, tc.neededColumns)
			if res := splitter.CanSplitSpanIntoFamilySpans(tc.prefixLen, tc.containsNull); res != tc.canSplit {
				t.Errorf("expected result to be %v, but found %v", tc.canSplit, res)
			}
		})
	}
}
