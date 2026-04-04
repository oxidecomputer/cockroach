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

// CreateDropOwnedByCounter returns a counter to increment for the DROP OWNED BY command.
func CreateDropOwnedByCounter() telemetry.Counter {
	return telemetry.GetCounter("sql.drop_owned_by")
}
