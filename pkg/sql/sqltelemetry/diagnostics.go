// Copyright 2019 The Cockroach Authors.
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

import (
	"crypto/sha256"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
)

// StatementDiagnosticsCollectedCounter is to be incremented whenever a query is
// run with diagnostic collection (as a result of a user request through the
// UI). This does not include diagnostics collected through
// EXPLAIN ANALYZE (DEBUG), which has a separate counter.
// distributed across multiple nodes.
var StatementDiagnosticsCollectedCounter = telemetry.GetCounterOnce("sql.diagnostics.collected")

// HashedFeatureCounter returns a counter for the specified feature which hashes
// the feature name before reporting. This allows us to have a built-in which
// reports counts arbitrary feature names without risking its being used to
// transmit sensitive data, since only known hashes will be meaningful to
// the Cockroach Labs team.
func HashedFeatureCounter(feature string) telemetry.Counter {
	sum := sha256.Sum256([]byte(feature))
	return telemetry.GetCounter(fmt.Sprintf("sql.hashed.%x", sum))
}
