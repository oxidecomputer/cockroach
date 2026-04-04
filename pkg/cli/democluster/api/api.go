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

package democlusterapi

import (
	"context"
	"io"
)

// DemoCluster represents the subset of the API of a demo cluster
// that is exposed to the SQL shell. It only contains the part
// of the API that do not create a strong dependency on CockroachDB's
// server package and machinery.
type DemoCluster interface {
	// ListDemoNodes produces a listing of servers on the specified
	// writer. If justOne is specified, only the first node is listed.
	// Listing is printed to 'w'. Errors/warnings are printed to 'ew'.
	ListDemoNodes(w, ew io.Writer, justOne bool)

	// AddNode creates a new node with the given locality string.
	AddNode(ctx context.Context, localityString string) (newNodeID int32, err error)

	// GetLocality retrieves the locality of the given node.
	GetLocality(nodeID int32) string

	// NumNodes returns the number of nodes.
	NumNodes() int

	// DrainAndShutdown shuts down a node gracefully.
	DrainAndShutdown(ctx context.Context, nodeID int32) error

	// RestartNode starts the given node. The node must be down
	// prior to the call.
	RestartNode(ctx context.Context, nodeID int32) error

	// Decommission decommissions the given node.
	Decommission(ctx context.Context, nodeID int32) error

	// Recommission recommissions the given node.
	Recommission(ctx context.Context, nodeID int32) error
}
