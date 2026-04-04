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

package spanconfigkvsubscriber

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
)

// noopKVSubscriber is a KVSubscriber that no-ops and is always up-to-date.
// Intended for tests that do not make use of the span configurations
// infrastructure.
type noopKVSubscriber struct {
	clock *hlc.Clock
}

var _ spanconfig.KVSubscriber = &noopKVSubscriber{}

// NewNoopSubscriber returns a new no-op KVSubscriber.
func NewNoopSubscriber(clock *hlc.Clock) spanconfig.KVSubscriber {
	return &noopKVSubscriber{
		clock: clock,
	}
}

// Subscribe is part of the spanconfig.KVSubsriber interface.
func (n *noopKVSubscriber) Subscribe(func(context.Context, roachpb.Span)) {}

// LastUpdated is part of the spanconfig.KVSubscriber interface.
func (n *noopKVSubscriber) LastUpdated() hlc.Timestamp {
	return n.clock.Now()
}

// NeedsSplit is part of the spanconfig.KVSubscriber interface.
func (n *noopKVSubscriber) NeedsSplit(context.Context, roachpb.RKey, roachpb.RKey) bool {
	return false
}

// ComputeSplitKey is part of the spanconfig.KVSubscriber interface.
func (n *noopKVSubscriber) ComputeSplitKey(
	context.Context, roachpb.RKey, roachpb.RKey,
) roachpb.RKey {
	return roachpb.RKey{}
}

// GetSpanConfigForKey is part of the spanconfig.KVSubscriber interface.
func (n *noopKVSubscriber) GetSpanConfigForKey(
	context.Context, roachpb.RKey,
) (roachpb.SpanConfig, error) {
	return roachpb.SpanConfig{}, nil
}

// GetProtectionTimestamps is part of the spanconfig.KVSubscriber interface.
func (n *noopKVSubscriber) GetProtectionTimestamps(
	context.Context, roachpb.Span,
) ([]hlc.Timestamp, hlc.Timestamp, error) {
	return nil, n.LastUpdated(), nil
}
