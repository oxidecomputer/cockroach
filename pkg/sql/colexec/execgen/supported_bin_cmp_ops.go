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

package execgen

import (
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree/treebin"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree/treecmp"
)

// BinaryOpName is a mapping from all binary operators that are supported by
// the vectorized engine to their names.
var BinaryOpName = map[treebin.BinaryOperatorSymbol]string{
	treebin.Bitand:            "Bitand",
	treebin.Bitor:             "Bitor",
	treebin.Bitxor:            "Bitxor",
	treebin.Plus:              "Plus",
	treebin.Minus:             "Minus",
	treebin.Mult:              "Mult",
	treebin.Div:               "Div",
	treebin.FloorDiv:          "FloorDiv",
	treebin.Mod:               "Mod",
	treebin.Pow:               "Pow",
	treebin.Concat:            "Concat",
	treebin.LShift:            "LShift",
	treebin.RShift:            "RShift",
	treebin.JSONFetchVal:      "JSONFetchVal",
	treebin.JSONFetchText:     "JSONFetchText",
	treebin.JSONFetchValPath:  "JSONFetchValPath",
	treebin.JSONFetchTextPath: "JSONFetchTextPath",
}

// ComparisonOpName is a mapping from all comparison operators that are
// supported by the vectorized engine to their names.
var ComparisonOpName = map[treecmp.ComparisonOperatorSymbol]string{
	treecmp.EQ: "EQ",
	treecmp.NE: "NE",
	treecmp.LT: "LT",
	treecmp.LE: "LE",
	treecmp.GT: "GT",
	treecmp.GE: "GE",
}
