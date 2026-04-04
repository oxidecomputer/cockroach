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

package schemaexpr

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// RowIndexedVarContainer is used to evaluate expressions over various rows.
type RowIndexedVarContainer struct {
	CurSourceRow tree.Datums

	// Because the rows we have might not be permuted in the same way as the
	// original table, we need to store a mapping between them.

	Cols    []catalog.Column
	Mapping catalog.TableColMap
}

var _ tree.IndexedVarContainer = &RowIndexedVarContainer{}

// IndexedVarEval implements tree.IndexedVarContainer.
func (r *RowIndexedVarContainer) IndexedVarEval(
	idx int, ctx *tree.EvalContext,
) (tree.Datum, error) {
	rowIdx, ok := r.Mapping.Get(r.Cols[idx].GetID())
	if !ok {
		return tree.DNull, nil
	}
	return r.CurSourceRow[rowIdx], nil
}

// IndexedVarResolvedType implements tree.IndexedVarContainer.
func (*RowIndexedVarContainer) IndexedVarResolvedType(idx int) *types.T {
	panic("unsupported")
}

// IndexedVarNodeFormatter implements tree.IndexedVarContainer.
func (*RowIndexedVarContainer) IndexedVarNodeFormatter(idx int) tree.NodeFormatter {
	return nil
}

// CannotWriteToComputedColError constructs a write error for a computed column.
func CannotWriteToComputedColError(colName string) error {
	return pgerror.Newf(pgcode.ObjectNotInPrerequisiteState,
		"cannot write directly to computed column %q", tree.ErrNameString(colName))
}
