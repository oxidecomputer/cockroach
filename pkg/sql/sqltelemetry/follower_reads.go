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

package sqltelemetry

import "github.com/cockroachdb/cockroach/pkg/server/telemetry"

// FollowerReadDisabledCCLCounter is to be increment every time follower reads
// are requested but unavailable due to not having the CCL build.
var FollowerReadDisabledCCLCounter = telemetry.GetCounterOnce("follower_reads.disabled.ccl")

// FollowerReadDisabledNoEnterpriseLicense is to be incremented every time follower reads
// are requested but unavailable due to not having enterprise enabled.
var FollowerReadDisabledNoEnterpriseLicense = telemetry.GetCounterOnce("follower_reads.disabled.no_enterprise_license")
