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
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// buildLimit adds Limit and Offset operators according to the Limit clause.
//
// parentScope is the scope for the LIMIT/OFFSET expressions; this is not the
// same as inScope, because statements like:
//   SELECT k FROM kv LIMIT k
// are not valid.
func (b *Builder) buildLimit(limit *tree.Limit, parentScope, inScope *scope) {
	if limit.Offset != nil {
		input := inScope.expr
		offset := b.resolveAndBuildScalar(
			limit.Offset, types.Int, exprKindOffset, tree.RejectSpecial, parentScope,
		)
		inScope.expr = b.factory.ConstructOffset(input, offset, inScope.makeOrderingChoice())
	}
	if limit.Count != nil {
		input := inScope.expr
		limit := b.resolveAndBuildScalar(
			limit.Count, types.Int, exprKindLimit, tree.RejectSpecial, parentScope,
		)
		inScope.expr = b.factory.ConstructLimit(input, limit, inScope.makeOrderingChoice())
	}
}
