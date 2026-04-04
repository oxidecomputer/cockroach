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

package norm

import (
	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
)

// CanInlineWith returns whether or not it's valid to inline binding in expr.
// This is the case when materialize is explicitly set to false, or when:
// 1. binding has no volatile expressions (because once it's inlined, there's no
//    guarantee it will be executed fully), and
// 2. binding is referenced at most once in expr.
func (c *CustomFuncs) CanInlineWith(binding, expr memo.RelExpr, private *memo.WithPrivate) bool {
	// If materialization is set, ignore the checks below.
	if private.Mtr.Set {
		return !private.Mtr.Materialize
	}
	if binding.Relational().VolatilitySet.HasVolatile() {
		return false
	}
	return memo.WithUses(expr)[private.ID].Count <= 1
}

// InlineWith replaces all references to the With expression in input (via
// WithScans) with its definition.
func (c *CustomFuncs) InlineWith(binding, input memo.RelExpr, priv *memo.WithPrivate) memo.RelExpr {
	var replace ReplaceFunc
	replace = func(nd opt.Expr) opt.Expr {
		switch t := nd.(type) {
		case *memo.WithScanExpr:
			if t.With == priv.ID {
				// TODO(justin): it might be worth carefully walking the tree and
				// renaming variables as we do this replacement so that this projection
				// is unnecessary (assuming there's at most one reference to the
				// WithScan, which might be false if we heuristically inline multiple
				// times in the future).
				projections := make(memo.ProjectionsExpr, len(t.InCols))
				for i := range t.InCols {
					projections[i] = c.f.ConstructProjectionsItem(
						c.f.ConstructVariable(t.InCols[i]),
						t.OutCols[i],
					)
				}
				return c.f.ConstructProject(binding, projections, opt.ColSet{})
			}
			// TODO(justin): should apply joins block inlining because they can lead
			// to expressions being executed multiple times?
		}
		return c.f.Replace(nd, replace)
	}

	return replace(input).(memo.RelExpr)
}
