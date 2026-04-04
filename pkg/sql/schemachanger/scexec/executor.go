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

package scexec

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
)

// ExecuteStage executes the provided ops. The ops must all be of the same type.
func ExecuteStage(ctx context.Context, deps Dependencies, ops []scop.Op) error {
	// It is perfectly valid to have empty stage after optimizations /
	// transformations.
	if len(ops) == 0 {
		log.Infof(ctx, "skipping execution, no operations in this stage")
		return nil
	}
	typ := ops[0].Type()
	switch typ {
	case scop.MutationType:
		return executeDescriptorMutationOps(ctx, deps, ops)
	case scop.BackfillType:
		return executeBackfillOps(ctx, deps, ops)
	case scop.ValidationType:
		return executeValidationOps(ctx, deps, ops)
	default:
		return errors.AssertionFailedf("unknown ops type %d", typ)
	}
}
