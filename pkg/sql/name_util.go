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

package sql

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/clusterversion"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/catalogkeys"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func (p *planner) dropNamespaceEntry(
	ctx context.Context, b *kv.Batch, desc catalog.MutableDescriptor,
) {
	// Delete current namespace entry.
	deleteNamespaceEntryAndMaybeAddDrainingName(ctx, b, p, desc, desc)
}

func (p *planner) renameNamespaceEntry(
	ctx context.Context, b *kv.Batch, oldNameKey catalog.NameKey, desc catalog.MutableDescriptor,
) {
	// Delete old namespace entry.
	deleteNamespaceEntryAndMaybeAddDrainingName(ctx, b, p, oldNameKey, desc)

	// Write new namespace entry.
	marshalledKey := catalogkeys.EncodeNameKey(p.ExecCfg().Codec, desc)
	if p.extendedEvalCtx.Tracing.KVTracingEnabled() {
		log.VEventf(ctx, 2, "CPut %s -> %d", marshalledKey, desc.GetID())
	}
	b.CPut(marshalledKey, desc.GetID(), nil)
}

func deleteNamespaceEntryAndMaybeAddDrainingName(
	ctx context.Context,
	b *kv.Batch,
	p *planner,
	nameKeyToDelete catalog.NameKey,
	desc catalog.MutableDescriptor,
) {
	if !p.execCfg.Settings.Version.IsActive(ctx, clusterversion.AvoidDrainingNames) {
		desc.AddDrainingName(descpb.NameInfo{
			ParentID:       nameKeyToDelete.GetParentID(),
			ParentSchemaID: nameKeyToDelete.GetParentSchemaID(),
			Name:           nameKeyToDelete.GetName(),
		})
		return
	}
	marshalledKey := catalogkeys.EncodeNameKey(p.ExecCfg().Codec, nameKeyToDelete)
	if p.extendedEvalCtx.Tracing.KVTracingEnabled() {
		log.VEventf(ctx, 2, "Del %s", marshalledKey)
	}
	b.Del(marshalledKey)
}
