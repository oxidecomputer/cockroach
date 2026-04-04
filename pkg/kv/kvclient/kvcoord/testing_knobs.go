// Copyright 2018 The Cockroach Authors.
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

package kvcoord

import "github.com/cockroachdb/cockroach/pkg/base"

// ClientTestingKnobs contains testing options that dictate the behavior
// of the key-value client.
type ClientTestingKnobs struct {
	// The RPC dispatcher. Defaults to grpc but can be changed here for
	// testing purposes.
	TransportFactory TransportFactory

	// DontConsiderConnHealth, if set, makes the GRPCTransport not take into
	// consideration the connection health when deciding the ordering for
	// replicas. When not set, replicas on nodes with unhealthy connections are
	// deprioritized.
	DontConsiderConnHealth bool

	// The maximum number of times a txn will attempt to refresh its
	// spans for a single transactional batch.
	// 0 means use a default. -1 means disable refresh.
	MaxTxnRefreshAttempts int

	// CondenseRefreshSpansFilter, if set, is called when the span refresher is
	// considering condensing the refresh spans. If it returns false, condensing
	// will not be attempted and the span refresher will behave as if condensing
	// failed to save enough memory.
	CondenseRefreshSpansFilter func() bool

	// LatencyFunc, if set, overrides RPCContext.RemoteClocks.Latency as the
	// function used by the DistSender to order replicas for follower reads.
	LatencyFunc LatencyFunc

	// If set, the DistSender will try the replicas in the order they appear in
	// the descriptor, instead of trying to reorder them by latency. The knob
	// only applies to requests sent with the LEASEHOLDER routing policy.
	DontReorderReplicas bool

	// DisableCommitSanityCheck allows "setting" the DisableCommitSanityCheck to
	// true without actually overriding the variable.
	DisableCommitSanityCheck bool
}

var _ base.ModuleTestingKnobs = &ClientTestingKnobs{}

// ModuleTestingKnobs is part of the base.ModuleTestingKnobs interface.
func (*ClientTestingKnobs) ModuleTestingKnobs() {}
