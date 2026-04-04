// Copyright 2021 The Cockroach Authors.
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

package tree

import "strings"

// SchemaFeatureName feature name for a given statement, which can be used
// to detect via the feature check functions if the schema change is allowed.
type SchemaFeatureName string

// GetSchemaFeatureNameFromStmt takes a statement and converts it to a schema
// feature name, which can be enabled or disabled via a feature flag.
func GetSchemaFeatureNameFromStmt(stmt Statement) SchemaFeatureName {
	statementTag := stmt.StatementTag()
	statementInfo := strings.Split(statementTag, " ")
	// Only grab the first two words (i.e. ALTER TABLE, etc..).
	if len(statementInfo) >= 2 {
		return SchemaFeatureName(statementInfo[0] + " " + statementInfo[1])
	}
	return SchemaFeatureName(statementInfo[0])
}
