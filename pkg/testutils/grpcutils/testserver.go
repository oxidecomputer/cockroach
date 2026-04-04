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

package grpcutils

import (
	"context"

	"github.com/gogo/protobuf/types"
)

// TestServerImpl backs the Test service.
type TestServerImpl struct {
	UU func(context.Context, *types.Any) (*types.Any, error) // UnaryUnary
	US func(*types.Any, GRPCTest_UnaryStreamServer) error    // UnaryStream
	SU func(server GRPCTest_StreamUnaryServer) error         // StreamUnary
	SS func(server GRPCTest_StreamStreamServer) error        // StreamStream
}

var _ GRPCTestServer = (*TestServerImpl)(nil)

// UnaryUnary implements GRPCTestServer.
func (s *TestServerImpl) UnaryUnary(ctx context.Context, any *types.Any) (*types.Any, error) {
	return s.UU(ctx, any)
}

// UnaryStream implements GRPCTestServer.
func (s *TestServerImpl) UnaryStream(any *types.Any, server GRPCTest_UnaryStreamServer) error {
	return s.US(any, server)
}

// StreamUnary implements GRPCTestServer.
func (s *TestServerImpl) StreamUnary(server GRPCTest_StreamUnaryServer) error {
	return s.SU(server)
}

// StreamStream implements GRPCTestServer.
func (s *TestServerImpl) StreamStream(server GRPCTest_StreamStreamServer) error {
	return s.SS(server)
}
