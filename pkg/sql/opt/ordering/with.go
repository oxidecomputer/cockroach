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

package ordering

import (
	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/props"
)

func withCanProvideOrdering(expr memo.RelExpr, required *props.OrderingChoice) bool {
	// With operator can always pass through ordering to its main input.
	return true
}

func withBuildChildReqOrdering(
	parent memo.RelExpr, required *props.OrderingChoice, childIdx int,
) props.OrderingChoice {
	switch childIdx {
	case 0:
		return parent.(*memo.WithExpr).BindingOrdering

	case 1:
		// We can pass through any required ordering to the main query.
		return *required
	}
	return props.OrderingChoice{}
}

func withBuildProvided(expr memo.RelExpr, required *props.OrderingChoice) opt.Ordering {
	w := expr.(*memo.WithExpr)
	return w.Main.ProvidedPhysical().Ordering
}
