// Copyright 2015 The Cockroach Authors.
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

package sql

import (
	"context"
	"fmt"
	"math"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// limitNode represents a node that limits the number of rows
// returned or only return them past a given number (offset).
type limitNode struct {
	plan       planNode
	countExpr  tree.TypedExpr
	offsetExpr tree.TypedExpr
}

func (n *limitNode) startExec(params runParams) error {
	panic("limitNode cannot be run in local mode")
}

func (n *limitNode) Next(params runParams) (bool, error) {
	panic("limitNode cannot be run in local mode")
}

func (n *limitNode) Values() tree.Datums {
	panic("limitNode cannot be run in local mode")
}

func (n *limitNode) Close(ctx context.Context) {
	n.plan.Close(ctx)
}

// evalLimit evaluates the Count and Offset fields. If Count is missing, the
// value is MaxInt64. If Offset is missing, the value is 0
func evalLimit(
	evalCtx *tree.EvalContext, countExpr, offsetExpr tree.TypedExpr,
) (count, offset int64, err error) {
	count = math.MaxInt64
	offset = 0

	data := []struct {
		name string
		src  tree.TypedExpr
		dst  *int64
	}{
		{"LIMIT", countExpr, &count},
		{"OFFSET", offsetExpr, &offset},
	}

	for _, datum := range data {
		if datum.src != nil {
			dstDatum, err := datum.src.Eval(evalCtx)
			if err != nil {
				return count, offset, err
			}

			if dstDatum == tree.DNull {
				// Use the default value.
				continue
			}

			dstDInt := tree.MustBeDInt(dstDatum)
			val := int64(dstDInt)
			if val < 0 {
				return count, offset, fmt.Errorf("negative value for %s", datum.name)
			}
			*datum.dst = val
		}
	}
	return count, offset, nil
}
