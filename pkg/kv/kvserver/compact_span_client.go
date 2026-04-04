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

package kvserver

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/rpc"
	"github.com/cockroachdb/cockroach/pkg/rpc/nodedialer"
	"github.com/cockroachdb/errors"
)

// CompactEngineSpanClient is used to request compaction for a span
// of data on a store.
type CompactEngineSpanClient struct {
	nd *nodedialer.Dialer
}

// NewCompactEngineSpanClient constructs a new CompactEngineSpanClient.
func NewCompactEngineSpanClient(nd *nodedialer.Dialer) *CompactEngineSpanClient {
	return &CompactEngineSpanClient{nd: nd}
}

// CompactEngineSpan is a tree.CompactEngineSpanFunc.
func (c *CompactEngineSpanClient) CompactEngineSpan(
	ctx context.Context, nodeID, storeID int32, startKey, endKey []byte,
) error {
	conn, err := c.nd.Dial(ctx, roachpb.NodeID(nodeID), rpc.DefaultClass)
	if err != nil {
		return errors.Wrapf(err, "could not dial node ID %d", nodeID)
	}
	client := NewPerStoreClient(conn)
	req := &CompactEngineSpanRequest{
		StoreRequestHeader: StoreRequestHeader{
			NodeID:  roachpb.NodeID(nodeID),
			StoreID: roachpb.StoreID(storeID),
		},
		Span: roachpb.Span{Key: roachpb.Key(startKey), EndKey: roachpb.Key(endKey)},
	}
	_, err = client.CompactEngineSpan(ctx, req)
	return err
}
