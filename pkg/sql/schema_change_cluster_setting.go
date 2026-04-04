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

package sql

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/featureflag"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

// featureSchemaChangeEnabled is the cluster setting used to enable and disable
// any features that require schema changes. Documentation for which features
// are covered TBD.
var featureSchemaChangeEnabled = settings.RegisterBoolSetting(
	settings.TenantWritable,
	"feature.schema_change.enabled",
	"set to true to enable schema changes, false to disable; default is true",
	featureflag.FeatureFlagEnabledDefault,
).WithPublic()

// checkSchemaChangeEnabled is a method that wraps the featureflag.CheckEnabled
// method specifically for all features that are categorized as schema changes.
func checkSchemaChangeEnabled(
	ctx context.Context, execCfg *ExecutorConfig, schemaFeatureName string,
) error {
	if err := featureflag.CheckEnabled(
		ctx,
		execCfg,
		featureSchemaChangeEnabled,
		fmt.Sprintf("%s is part of the schema change category, which", schemaFeatureName),
	); err != nil {
		return err
	}
	return nil
}

// CheckFeature checks if a schema change feature is allowed.
func (p *planner) CheckFeature(ctx context.Context, featureName tree.SchemaFeatureName) error {
	return checkSchemaChangeEnabled(ctx, p.ExecCfg(), string(featureName))
}
