// Copyright 2022 The Cockroach Authors.
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

package contention

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/server/serverpb"
	"github.com/cockroachdb/cockroach/pkg/sql/contentionpb"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
)

type fakeStatusServer struct {
	data          map[uuid.UUID]roachpb.TransactionFingerprintID
	injectedError error
}

func newFakeStatusServer() *fakeStatusServer {
	return &fakeStatusServer{
		data:          make(map[uuid.UUID]roachpb.TransactionFingerprintID),
		injectedError: nil,
	}
}

func (f *fakeStatusServer) txnIDResolution(
	_ context.Context, req *serverpb.TxnIDResolutionRequest,
) (*serverpb.TxnIDResolutionResponse, error) {
	if f.injectedError != nil {
		return nil, f.injectedError
	}

	resp := &serverpb.TxnIDResolutionResponse{
		ResolvedTxnIDs: make([]contentionpb.ResolvedTxnID, 0),
	}

	for _, txnID := range req.TxnIDs {
		if txnFingerprintID, ok := f.data[txnID]; ok {
			resp.ResolvedTxnIDs = append(resp.ResolvedTxnIDs, contentionpb.ResolvedTxnID{
				TxnID:            txnID,
				TxnFingerprintID: txnFingerprintID,
			})
		}
	}

	return resp, nil
}

func (f *fakeStatusServer) setTxnIDEntry(
	txnID uuid.UUID, txnFingerprintID roachpb.TransactionFingerprintID,
) {
	f.data[txnID] = txnFingerprintID
}

type fakeStatusServerCluster map[string]*fakeStatusServer

func newFakeStatusServerCluster() fakeStatusServerCluster {
	return make(fakeStatusServerCluster)
}

func (f fakeStatusServerCluster) getStatusServer(coordinatorID string) *fakeStatusServer {
	statusServer, ok := f[coordinatorID]
	if !ok {
		statusServer = newFakeStatusServer()
		f[coordinatorID] = statusServer
	}
	return statusServer
}

func (f fakeStatusServerCluster) txnIDResolution(
	ctx context.Context, req *serverpb.TxnIDResolutionRequest,
) (*serverpb.TxnIDResolutionResponse, error) {
	return f.getStatusServer(req.CoordinatorID).txnIDResolution(ctx, req)
}

func (f fakeStatusServerCluster) setTxnIDEntry(
	coordinatorNodeID string, txnID uuid.UUID, txnFingerprintID roachpb.TransactionFingerprintID,
) {
	f.getStatusServer(coordinatorNodeID).setTxnIDEntry(txnID, txnFingerprintID)
}

func (f fakeStatusServerCluster) setStatusServerError(coordinatorNodeID string, err error) {
	f.getStatusServer(coordinatorNodeID).injectedError = err
}

func (f fakeStatusServerCluster) clear() {
	for k := range f {
		delete(f, k)
	}
}

func (f fakeStatusServerCluster) clearErrors() {
	for _, statusServer := range f {
		statusServer.injectedError = nil
	}
}
