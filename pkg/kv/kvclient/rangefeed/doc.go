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

// Package rangefeed provides a useful client abstraction atop of the rangefeed
// functionality exported by the DistSender.
//
// In particular, the abstraction exported by this package hooks up a stopper,
// and deals with retries upon errors, tracking resolved timestamps along the
// way.
package rangefeed

// TODO(ajwerner): Rework this logic to encapsulate the multi-span logic in
// changefeedccl/kvfeed. That code also deals with some schema interactions but
// it should be split into two layers. The primary limitation missing here is
// just the ability to watch multiple spans however the way that the KV feed
// manages internal state and sometimes triggers re-scanning would require some
// interface changes.
