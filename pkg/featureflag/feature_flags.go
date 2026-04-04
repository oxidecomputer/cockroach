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

package featureflag

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
)

// FeatureFlagEnabledDefault is used for the default value of all feature flag
// cluster settings.
const FeatureFlagEnabledDefault = true

// CheckEnabled checks whether a specific feature has been enabled or disabled
// via its cluster settings, and returns an error if it has been disabled. It
// increments a telemetry counter for any feature flag that is denied and logs
// the feature and category that was denied.
func CheckEnabled(
	ctx context.Context, config Config, s *settings.BoolSetting, featureName string,
) error {
	sv := config.SV()

	if enabled := s.Get(sv); !enabled {
		// Increment telemetry counter for feature flag denial.
		telemetry.Inc(sqltelemetry.FeatureDeniedByFeatureFlagCounter)

		// Increment metrics counter for feature flag denial.
		if config.GetFeatureFlagMetrics() == nil {
			log.Warningf(
				ctx,
				"executorConfig.FeatureFlagMetrics is uninitiated; feature %s was attempted but disabled via cluster settings",
				featureName,
			)
			return pgerror.Newf(
				pgcode.OperatorIntervention,
				"feature %s was disabled by the database administrator",
				featureName,
			)
		}
		config.GetFeatureFlagMetrics().FeatureDenialMetric.Inc(1)

		if log.V(2) {
			log.Warningf(
				ctx,
				"feature %s was attempted but disabled via cluster settings",
				featureName,
			)
		}

		return pgerror.Newf(
			pgcode.OperatorIntervention,
			"feature %s was disabled by the database administrator",
			featureName,
		)
	}
	return nil
}

// metaFeatureDenialMetric is a metric counting the statements denied by a
//feature flag.
var metaFeatureDenialMetric = metric.Metadata{
	Name:        "sql.feature_flag_denial",
	Help:        "Counter of the number of statements denied by a feature flag",
	Measurement: "Statements",
	Unit:        metric.Unit_COUNT,
}

// DenialMetrics is a struct corresponding to any metrics related to feature
// flag denials. Future metrics related to feature flags should be added to
// this struct.
type DenialMetrics struct {
	FeatureDenialMetric *metric.Counter
}

// MetricStruct makes FeatureFlagMetrics a metric.Struct.
func (s *DenialMetrics) MetricStruct() {}

var _ metric.Struct = (*DenialMetrics)(nil)

// NewFeatureFlagMetrics constructs a new FeatureFlagMetrics metric.Struct.
func NewFeatureFlagMetrics() *DenialMetrics {
	return &DenialMetrics{
		FeatureDenialMetric: metric.NewCounter(metaFeatureDenialMetric),
	}
}

// Config is an interface used to pass required values from ExecutorConfig to
// check if a feature flag is enabled.
type Config interface {
	SV() *settings.Values
	GetFeatureFlagMetrics() *DenialMetrics
}
