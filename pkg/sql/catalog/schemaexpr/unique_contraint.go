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

package schemaexpr

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
)

// ValidateUniqueWithoutIndexPredicate verifies that an expression is a valid
// unique without index predicate. If the expression is valid, it returns the
// serialized expression with the columns dequalified.
//
// A predicate expression is valid if all of the following are true:
//
//   - It results in a boolean.
//   - It refers only to columns in the table.
//   - It does not include subqueries.
//   - It does not include non-immutable, aggregate, window, or set returning
//     functions.
//
func ValidateUniqueWithoutIndexPredicate(
	ctx context.Context,
	tn tree.TableName,
	desc catalog.TableDescriptor,
	pred tree.Expr,
	semaCtx *tree.SemaContext,
) (string, error) {
	expr, _, _, err := DequalifyAndValidateExpr(
		ctx,
		desc,
		pred,
		types.Bool,
		"unique without index predicate",
		semaCtx,
		tree.VolatilityImmutable,
		&tn,
	)
	if err != nil {
		return "", err
	}
	return expr, nil
}
