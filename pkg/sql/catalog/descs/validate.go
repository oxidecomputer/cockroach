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

package descs

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catconstants"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/internal/catkv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/internal/validate"
)

// Validate returns any descriptor validation errors after validating using the
// descriptor collection for retrieving referenced descriptors and namespace
// entries, if applicable.
func (tc *Collection) Validate(
	ctx context.Context,
	txn *kv.Txn,
	telemetry catalog.ValidationTelemetry,
	targetLevel catalog.ValidationLevel,
	descriptors ...catalog.Descriptor,
) (err error) {
	vd := tc.newValidationDereferencer(txn)
	version := tc.settings.Version.ActiveVersion(ctx)
	return validate.Validate(
		ctx,
		version,
		vd,
		telemetry,
		targetLevel,
		descriptors...).CombinedError()
}

// ValidateUncommittedDescriptors validates all uncommitted descriptors.
// Validation includes cross-reference checks. Referenced descriptors are
// read from the store unless they happen to also be part of the uncommitted
// descriptor set. We purposefully avoid using leased descriptors as those may
// be one version behind, in which case it's possible (and legitimate) that
// those are missing back-references which would cause validation to fail.
func (tc *Collection) ValidateUncommittedDescriptors(ctx context.Context, txn *kv.Txn) (err error) {
	if tc.skipValidationOnWrite || !ValidateOnWriteEnabled.Get(&tc.settings.SV) {
		return nil
	}
	descs := tc.uncommitted.getUncommittedDescriptorsForValidation()
	if len(descs) == 0 {
		return nil
	}
	return tc.Validate(ctx, txn, catalog.ValidationWriteTelemetry, catalog.ValidationLevelAllPreTxnCommit, descs...)
}

func (tc *Collection) newValidationDereferencer(txn *kv.Txn) validate.ValidationDereferencer {
	return &collectionBackedDereferencer{tc: tc, txn: txn}
}

// collectionBackedDereferencer wraps a Collection to implement the
// validate.ValidationDereferencer interface for validation.
type collectionBackedDereferencer struct {
	tc  *Collection
	txn *kv.Txn
}

var _ validate.ValidationDereferencer = &collectionBackedDereferencer{}

// DereferenceDescriptors implements the validate.ValidationDereferencer
// interface by leveraging the collection's uncommitted descriptors.
func (c collectionBackedDereferencer) DereferenceDescriptors(
	ctx context.Context, version clusterversion.ClusterVersion, reqs []descpb.ID,
) (ret []catalog.Descriptor, _ error) {
	ret = make([]catalog.Descriptor, len(reqs))
	fallbackReqs := make([]descpb.ID, 0, len(reqs))
	fallbackRetIndexes := make([]int, 0, len(reqs))
	for i, id := range reqs {
		desc, err := c.fastDescLookup(ctx, id)
		if err != nil {
			return nil, err
		}
		if desc == nil {
			fallbackReqs = append(fallbackReqs, id)
			fallbackRetIndexes = append(fallbackRetIndexes, i)
		} else {
			ret[i] = desc
		}
	}
	if len(fallbackReqs) > 0 {
		fallbackRet, err := catkv.GetCrossReferencedDescriptorsForValidation(
			ctx,
			version,
			c.tc.codec(),
			c.txn,
			fallbackReqs,
		)
		if err != nil {
			return nil, err
		}
		for j, desc := range fallbackRet {
			if desc == nil {
				continue
			}
			if uc, _ := c.tc.uncommitted.getImmutableByID(desc.GetID()); uc == nil {
				desc, err = c.tc.uncommitted.add(desc.NewBuilder().BuildExistingMutable(), notValidatedYet)
				if err != nil {
					return nil, err
				}
			}
			ret[fallbackRetIndexes[j]] = desc
		}
	}
	return ret, nil
}

func (c collectionBackedDereferencer) fastDescLookup(
	ctx context.Context, id descpb.ID,
) (catalog.Descriptor, error) {
	if uc, _ := c.tc.uncommitted.getImmutableByID(id); uc != nil {
		return uc, nil
	}
	return nil, nil
}

// DereferenceDescriptorIDs implements the validate.ValidationDereferencer
// interface.
func (c collectionBackedDereferencer) DereferenceDescriptorIDs(
	ctx context.Context, reqs []descpb.NameInfo,
) (ret []descpb.ID, _ error) {
	ret = make([]descpb.ID, len(reqs))
	fallbackReqs := make([]descpb.NameInfo, 0, len(reqs))
	fallbackRetIndexes := make([]int, 0, len(reqs))
	for i, req := range reqs {
		found, id, err := c.fastNamespaceLookup(ctx, req)
		if err != nil {
			return nil, err
		}
		if found {
			ret[i] = id
		} else {
			fallbackReqs = append(fallbackReqs, req)
			fallbackRetIndexes = append(fallbackRetIndexes, i)
		}
	}
	if len(fallbackReqs) > 0 {
		// TODO(postamar): actually use the Collection here instead,
		// either by calling the Collection's methods or by caching the results
		// of this call in the Collection.
		fallbackRet, err := catkv.LookupIDs(ctx, c.txn, c.tc.codec(), fallbackReqs)
		if err != nil {
			return nil, err
		}
		for j, id := range fallbackRet {
			ret[fallbackRetIndexes[j]] = id
		}
	}
	return ret, nil
}

func (c collectionBackedDereferencer) fastNamespaceLookup(
	ctx context.Context, req descpb.NameInfo,
) (found bool, id descpb.ID, err error) {
	// Handle special cases.
	// TODO(postamar): namespace lookups should go through Collection
	switch req.ParentID {
	case descpb.InvalidID:
		if req.ParentSchemaID == descpb.InvalidID && req.Name == catconstants.SystemDatabaseName {
			// Looking up system database ID, which is hard-coded.
			return true, keys.SystemDatabaseID, nil
		}
	case keys.SystemDatabaseID:
		// Looking up system database objects, which are cached.
		id = c.tc.kv.systemNamespace.lookup(req.ParentSchemaID, req.Name)
		return id != descpb.InvalidID, id, nil
	}
	return false, descpb.InvalidID, nil
}
