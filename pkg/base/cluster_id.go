// Copyright 2017 The Cockroach Authors.
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

package base

import (
	"context"
	"sync/atomic"

	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
)

// ClusterIDContainer is used to share a single Cluster ID instance between
// multiple layers. It allows setting and getting the value. Once a value is
// set, the value cannot change.
//
// The cluster ID is determined on startup as follows:
// - If there are existing stores, their cluster ID is used.
// - If the node is bootstrapping, a new UUID is generated.
// - Otherwise, it is determined via gossip with other nodes.
type ClusterIDContainer struct {
	clusterID atomic.Value // uuid.UUID
	// OnSet, if non-nil, is called after the ID is set with the new value.
	OnSet func(uuid.UUID)
}

// String returns the cluster ID, or "?" if it is unset.
func (c *ClusterIDContainer) String() string {
	val := c.Get()
	if val == uuid.Nil {
		return "?"
	}
	return val.String()
}

// Get returns the current cluster ID; uuid.Nil if it is unset.
func (c *ClusterIDContainer) Get() uuid.UUID {
	v := c.clusterID.Load()
	if v == nil {
		return uuid.Nil
	}
	return v.(uuid.UUID)
}

// Set sets the current cluster ID. If it is already set, the value must match.
func (c *ClusterIDContainer) Set(ctx context.Context, val uuid.UUID) {
	// NOTE: this compare-and-swap is intentionally racy and won't catch all
	// cases where two different cluster IDs are set. That's ok, as this is
	// just a sanity check. But if we decide to care, we can use the new
	// (*atomic.Value).CompareAndSwap API introduced in go1.17.
	cur := c.Get()
	if cur == uuid.Nil {
		c.clusterID.Store(val)
		if log.V(2) {
			log.Infof(ctx, "ClusterID set to %s", val)
		}
	} else if cur != val {
		log.Fatalf(ctx, "different ClusterIDs set: %s, then %s", cur, val)
	}
	if c.OnSet != nil {
		c.OnSet(val)
	}
}

// Reset changes the ClusterID regardless of the old value.
//
// Should only be used in testing code.
func (c *ClusterIDContainer) Reset(val uuid.UUID) {
	c.clusterID.Store(val)
}
