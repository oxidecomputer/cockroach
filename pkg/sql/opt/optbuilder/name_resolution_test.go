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

package optbuilder

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/colinfo/colinfotestutils"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

var _ colinfotestutils.ColumnItemResolverTester = &scope{}

// GetColumnItemResolver is part of the sqlutils.ColumnItemResolverTester
// interface.
func (s *scope) GetColumnItemResolver() colinfo.ColumnItemResolver {
	return s
}

// AddTable is part of the sqlutils.ColumnItemResolverTester interface.
func (s *scope) AddTable(tabName tree.TableName, colNames []tree.Name) {
	for _, col := range colNames {
		s.cols = append(s.cols, scopeColumn{name: scopeColName(col), table: tabName})
	}
}

// ResolveQualifiedStarTestResults is part of the
// sqlutils.ColumnItemResolverTester interface.
func (s *scope) ResolveQualifiedStarTestResults(
	srcName *tree.TableName, srcMeta colinfo.ColumnSourceMeta,
) (string, string, error) {
	s, ok := srcMeta.(*scope)
	if !ok {
		return "", "", fmt.Errorf("resolver did not return *scope, found %T instead", srcMeta)
	}
	nl := make(tree.NameList, 0, len(s.cols))
	for i := range s.cols {
		col := s.cols[i]
		if col.table == *srcName && col.visibility == visible {
			nl = append(nl, col.name.ReferenceName())
		}
	}
	return srcName.String(), nl.String(), nil
}

// ResolveColumnItemTestResults is part of the
// sqlutils.ColumnItemResolverTester interface.
func (s *scope) ResolveColumnItemTestResults(
	colRes colinfo.ColumnResolutionResult,
) (string, error) {
	col, ok := colRes.(*scopeColumn)
	if !ok {
		return "", fmt.Errorf("resolver did not return *scopeColumn, found %T instead", colRes)
	}
	return fmt.Sprintf("%s.%s", col.table.String(), col.name.ReferenceName()), nil
}

func TestResolveQualifiedStar(t *testing.T) {
	s := &scope{}
	colinfotestutils.RunResolveQualifiedStarTest(t, s)
}

func TestResolveColumnItem(t *testing.T) {
	s := &scope{}
	colinfotestutils.RunResolveColumnItemTest(t, s)
}
