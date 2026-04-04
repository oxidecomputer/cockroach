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

package stmtdiagnostics

import (
	"context"
	"time"
)

// InsertRequestInternal exposes the form of insert which returns the request ID
// as an int64 to tests in this package.
func (r *Registry) InsertRequestInternal(
	ctx context.Context, fprint string, minExecutionLatency time.Duration, expiresAfter time.Duration,
) (int64, error) {
	id, err := r.insertRequestInternal(ctx, fprint, minExecutionLatency, expiresAfter)
	return int64(id), err
}

// PollingInterval is exposed to override in tests.
var PollingInterval = pollingInterval
