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

package indexrec

import "testing"

func TestBuildOptAndHypTableMaps(t *testing.T) {
	tables, indexCols := testTablesAndIndexCols()
	table1 := tables[0]
	table2 := tables[1]
	indexCandidates := testIndexCandidates1(tables, indexCols)

	oldTables, hypTables := BuildOptAndHypTableMaps(indexCandidates)

	if oldTables[table1.ID()] != table1 {
		t.Errorf("expected table1 to be %+v,\n got %+v\n", table1, oldTables[table1.ID()])
	}

	if oldTables[table2.ID()] != table2 {
		t.Errorf("expected table2 to be %+v,\n got %+v\n", table2, oldTables[table2.ID()])
	}

	// A hypothetical table's index count is equivalent to its number of index
	// candidates plus the number of existing indexes.
	indexCountTable1 := len(indexCandidates[table1]) + table1.IndexCount()
	indexCountTable2 := len(indexCandidates[table2]) + table2.IndexCount()

	if hypTables[1].IndexCount() != indexCountTable1 {
		t.Errorf(
			"expected table1's index count to be %d, got %d\n",
			hypTables[1].IndexCount(),
			indexCountTable1,
		)
	}

	if hypTables[2].IndexCount() != indexCountTable2 {
		t.Errorf("expected table2's index count to be %d, got %d\n",
			hypTables[2].IndexCount(),
			indexCountTable2,
		)
	}
}
