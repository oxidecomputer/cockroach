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

import "github.com/cockroachdb/cockroach/pkg/sql/sem/tree"

// MakeHashShardComputeExpr creates the serialized computed expression for a hash shard
// column based on the column names and the number of buckets. The expression will be
// of the form:
//
//    mod(fnv32(crdb_internal.datums_to_bytes(...)),buckets)
//
func MakeHashShardComputeExpr(colNames []string, buckets int) *string {
	unresolvedFunc := func(funcName string) tree.ResolvableFunctionReference {
		return tree.ResolvableFunctionReference{
			FunctionReference: &tree.UnresolvedName{
				NumParts: 1,
				Parts:    tree.NameParts{funcName},
			},
		}
	}
	columnItems := func() tree.Exprs {
		exprs := make(tree.Exprs, len(colNames))
		for i := range exprs {
			exprs[i] = &tree.ColumnItem{ColumnName: tree.Name(colNames[i])}
		}
		return exprs
	}
	hashedColumnsExpr := func() tree.Expr {
		return &tree.FuncExpr{
			Func: unresolvedFunc("fnv32"),
			Exprs: tree.Exprs{
				&tree.FuncExpr{
					Func:  unresolvedFunc("crdb_internal.datums_to_bytes"),
					Exprs: columnItems(),
				},
			},
		}
	}
	modBuckets := func(expr tree.Expr) tree.Expr {
		return &tree.FuncExpr{
			Func: unresolvedFunc("mod"),
			Exprs: tree.Exprs{
				expr,
				tree.NewDInt(tree.DInt(buckets)),
			},
		}
	}
	res := tree.Serialize(modBuckets(hashedColumnsExpr()))
	return &res
}
