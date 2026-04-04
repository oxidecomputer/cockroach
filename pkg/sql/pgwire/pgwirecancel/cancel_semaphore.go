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

package pgwirecancel

import "github.com/cockroachdb/cockroach/pkg/util/quotapool"

// CancelSemaphore is a semaphore that limits the number of concurrent
// calls to the pgwire query cancellation endpoint. This is needed to avoid the
// risk of a DoS attack by malicious users that attempts to cancel random
// queries by spamming the request.
//
// We hard-code a limit of 256 concurrent pgwire cancel requests (per node).
// We also add a 1-second penalty for failed cancellation requests, meaning
// that an attacker needs 1 second per guess. With an attacker randomly
// guessing a 32-bit secret, it would take 2^24 seconds to hit one query. If
// we suppose there are 256 concurrent queries actively running on a node,
// then it would take 2^16 seconds (18 hours) to hit any one of them.
var CancelSemaphore = quotapool.NewIntPool("pgwire-cancel", 256)
