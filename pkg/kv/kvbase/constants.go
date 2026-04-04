// Copyright 2020 The Cockroach Authors.
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

package kvbase

// FollowerReadServingMsg is a log message that needs to be used for tests in
// other packages.
const FollowerReadServingMsg = "serving via follower read"

// RoutingRequestLocallyMsg is a log message that needs to be used for tests in
// other packages.
const RoutingRequestLocallyMsg = "sending request to local client"

// SpawningHeartbeatLoopMsg is a log message that needs to be used for tests in
// other packages.
const SpawningHeartbeatLoopMsg = "coordinator spawns heartbeat loop"
