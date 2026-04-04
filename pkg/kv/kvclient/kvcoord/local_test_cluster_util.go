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

package kvcoord

import (
	"context"
	"time"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/gossip"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/rpc"
	"github.com/cockroachdb/cockroach/pkg/rpc/nodedialer"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
)

// localTestClusterTransport augments senderTransport with an optional
// delay for each RPC, to simulate latency for benchmarking.
// TODO(bdarnell): there's probably a better place to put this.
type localTestClusterTransport struct {
	Transport
	latency time.Duration
}

func (l *localTestClusterTransport) SendNext(
	ctx context.Context, ba roachpb.BatchRequest,
) (*roachpb.BatchResponse, error) {
	if l.latency > 0 {
		time.Sleep(l.latency)
	}
	return l.Transport.SendNext(ctx, ba)
}

// InitFactoryForLocalTestCluster initializes a TxnCoordSenderFactory
// that can be used with LocalTestCluster.
func InitFactoryForLocalTestCluster(
	ctx context.Context,
	st *cluster.Settings,
	nodeDesc *roachpb.NodeDescriptor,
	tracer *tracing.Tracer,
	clock *hlc.Clock,
	latency time.Duration,
	stores kv.Sender,
	stopper *stop.Stopper,
	gossip *gossip.Gossip,
) kv.TxnSenderFactory {
	return NewTxnCoordSenderFactory(
		TxnCoordSenderFactoryConfig{
			AmbientCtx: log.MakeTestingAmbientContext(tracer),
			Settings:   st,
			Clock:      clock,
			Stopper:    stopper,
		},
		NewDistSenderForLocalTestCluster(ctx, st, nodeDesc, tracer, clock, latency, stores, stopper, gossip),
	)
}

// NewDistSenderForLocalTestCluster creates a DistSender for a LocalTestCluster.
func NewDistSenderForLocalTestCluster(
	ctx context.Context,
	st *cluster.Settings,
	nodeDesc *roachpb.NodeDescriptor,
	tracer *tracing.Tracer,
	clock *hlc.Clock,
	latency time.Duration,
	stores kv.Sender,
	stopper *stop.Stopper,
	g *gossip.Gossip,
) *DistSender {
	retryOpts := base.DefaultRetryOptions()
	retryOpts.Closer = stopper.ShouldQuiesce()
	rpcContext := rpc.NewInsecureTestingContext(ctx, clock, stopper)
	senderTransportFactory := SenderTransportFactory(tracer, stores)
	return NewDistSender(DistSenderConfig{
		AmbientCtx:         log.MakeTestingAmbientContext(tracer),
		Settings:           st,
		Clock:              clock,
		NodeDescs:          g,
		RPCContext:         rpcContext,
		RPCRetryOptions:    &retryOpts,
		nodeDescriptor:     nodeDesc,
		NodeDialer:         nodedialer.New(rpcContext, gossip.AddressResolver(g)),
		FirstRangeProvider: g,
		TestingKnobs: ClientTestingKnobs{
			TransportFactory: func(
				opts SendOptions,
				nodeDialer *nodedialer.Dialer,
				replicas ReplicaSlice,
			) (Transport, error) {
				transport, err := senderTransportFactory(opts, nodeDialer, replicas)
				if err != nil {
					return nil, err
				}
				return &localTestClusterTransport{transport, latency}, nil
			},
		},
	})
}
