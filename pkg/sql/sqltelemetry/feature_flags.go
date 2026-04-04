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

// FeatureDeniedByFeatureFlagCounter is a counter that is incremented every time a feature is
// denied via the feature flag cluster setting, for example. feature.schema_change.enabled = FALSE.
var FeatureDeniedByFeatureFlagCounter = telemetry.GetCounterOnce("sql.feature_flag_denied")
