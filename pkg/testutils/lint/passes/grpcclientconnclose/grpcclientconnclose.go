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

package grpcclientconnclose

import "github.com/cockroachdb/cockroach/pkg/testutils/lint/passes/forbiddenmethod"

// Analyzer checks for calls to (*grpc.ClientConn).Close. We mostly pull these
// objects from *rpc.Context, which manages their lifecycle.
// Errant calls to Close() disrupt the connection for all users.
// (Exported from forbiddenmethod.)
var Analyzer = forbiddenmethod.GRPCClientConnCloseAnalyzer
