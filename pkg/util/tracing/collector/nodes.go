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

package collector

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

// nodesFromNodeLiveness returns the IDs for all nodes that are currently part
// of the cluster (i.e. they haven't been decommissioned away). This list might
// also include nodes that are dead, in which case the RPC to collect traces
// from the dead node will timeout, and we will be able to better surface that
// error.
//
// It's important to note that this makes no guarantees about new nodes being
// added to the cluster. It's entirely possible for that to happen concurrently
// with the retrieval of the current set of nodes.
func nodesFromNodeLiveness(ctx context.Context, nl NodeLiveness) ([]roachpb.NodeID, error) {
	var ns []roachpb.NodeID
	ls, err := nl.GetLivenessesFromKV(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range ls {
		if l.Membership.Decommissioned() {
			continue
		}
		ns = append(ns, l.NodeID)
	}
	return ns, nil
}
