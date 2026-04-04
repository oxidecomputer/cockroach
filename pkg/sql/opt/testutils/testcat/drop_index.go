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

package testcat

import (
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

// DropIndex is a partial implementation of the DROP INDEX statement.
//
// It only supports dropping a secondary index with an unqualified name.
func (tc *Catalog) DropIndex(stmt *tree.DropIndex) {
	for _, tableIndexName := range stmt.IndexList {
		indexName := tableIndexName.Index.String()

		var foundTab *Table
		var idxOrd int
		for _, tab := range tc.Tables() {
			if idx, ok := findIndex(tab, indexName); ok {
				if foundTab != nil {
					panic(errors.Newf(
						`index name "%s" is ambiguous; dropping ambiguous indexes is not supported in the test catalog`,
						indexName,
					))
				}
				foundTab = tab
				idxOrd = idx.ordinal
			}
		}

		if foundTab == nil {
			panic(errors.Newf(`index "%s" does not exist`, indexName))
		}

		if idxOrd == 0 {
			panic(errors.Newf("dropping primary indexes is not supported in the test catalog"))
		}

		// Delete the index from the table.
		numIndexes := len(foundTab.Indexes)
		foundTab.Indexes[idxOrd] = foundTab.Indexes[numIndexes-1]
		foundTab.Indexes[idxOrd].ordinal = idxOrd
		foundTab.Indexes = foundTab.Indexes[:numIndexes-1]
	}
}

// findIndex returns the first index within tab that has an IdxName equal to
// name. If an index is found it returns the index and true, and nil and false
// otherwise.
func findIndex(tab *Table, name string) (*Index, bool) {
	for _, idx := range tab.Indexes {
		if idx.IdxName == name {
			return idx, true
		}
	}
	return nil, false
}
