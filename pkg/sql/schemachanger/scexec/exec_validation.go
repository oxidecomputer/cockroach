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

	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/errors"
)

func executeValidateUniqueIndex(
	ctx context.Context, deps Dependencies, op *scop.ValidateUniqueIndex,
) error {
	descs, err := deps.Catalog().MustReadImmutableDescriptors(ctx, op.TableID)
	if err != nil {
		return err
	}
	desc := descs[0]
	table, ok := desc.(catalog.TableDescriptor)
	if !ok {
		return catalog.WrapTableDescRefErr(desc.GetID(), catalog.NewDescriptorTypeError(desc))
	}
	index, err := table.FindIndexWithID(op.IndexID)
	if err != nil {
		return err
	}
	// Execute the validation operation as a root user.
	execOverride := sessiondata.InternalExecutorOverride{
		User: security.RootUserName(),
	}
	if index.GetType() == descpb.IndexDescriptor_FORWARD {
		err = deps.IndexValidator().ValidateForwardIndexes(ctx, table, []catalog.Index{index}, execOverride)
	} else {
		err = deps.IndexValidator().ValidateInvertedIndexes(ctx, table, []catalog.Index{index}, execOverride)
	}
	return err
}

func executeValidateCheckConstraint(
	ctx context.Context, deps Dependencies, op *scop.ValidateCheckConstraint,
) error {
	return errors.Errorf("executeValidateCheckConstraint is not implemented")
}

func executeValidationOps(ctx context.Context, deps Dependencies, execute []scop.Op) error {
	for _, op := range execute {
		switch op := op.(type) {
		case *scop.ValidateUniqueIndex:
			return executeValidateUniqueIndex(ctx, deps, op)
		case *scop.ValidateCheckConstraint:
			return executeValidateCheckConstraint(ctx, deps, op)
		default:
			panic("unimplemented")
		}
	}
	return nil
}
