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

package constraint

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/testutils/testcat"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

func TestColumns_RemapColumns(t *testing.T) {
	var md opt.Metadata
	catalog := testcat.New()
	_, err := catalog.ExecuteDDL("CREATE TABLE tab (a INT PRIMARY KEY, b INT, c INT, d INT);")
	if err != nil {
		t.Fatal(err)
	}
	tn := tree.NewUnqualifiedTableName("tab")
	tab := catalog.Table(tn)

	from := md.AddTable(tab, &tree.TableName{})
	to := md.AddTable(tab, &tree.TableName{})

	var originalColumns Columns
	originalColumns.Init([]opt.OrderingColumn{
		opt.MakeOrderingColumn(from.ColumnID(0), false /* descending */),
		opt.MakeOrderingColumn(from.ColumnID(2), true /* descending */),
		opt.MakeOrderingColumn(from.ColumnID(3), false /* descending */),
	})

	remappedColumns := originalColumns.RemapColumns(from, to)

	expected := "/1/-3/4"
	if originalColumns.String() != expected {
		t.Errorf("\noriginal Columns were changed: %s", originalColumns.String())
	}

	expected = "/7/-9/10"
	if remappedColumns.String() != expected {
		t.Errorf("\nexpected: %s\nactual: %s\n", expected, remappedColumns.String())
	}
}
