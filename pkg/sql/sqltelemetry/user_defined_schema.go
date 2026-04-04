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

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
)

// UserDefinedSchemaTelemetryType represents a type of user defined schema
// related operation to record telemetry for.
type UserDefinedSchemaTelemetryType int

const (
	_ UserDefinedSchemaTelemetryType = iota
	// UserDefinedSchemaCreate represents a CREATE SCHEMA command.
	UserDefinedSchemaCreate
	// UserDefinedSchemaAlter represents an ALTER SCHEMA command.
	UserDefinedSchemaAlter
	// UserDefinedSchemaDrop represents a DROP SCHEMA command.
	UserDefinedSchemaDrop
	// UserDefinedSchemaReparentDatabase represents an ALTER DATABASE ... CONVERT TO SCHEMA command.
	UserDefinedSchemaReparentDatabase
	// UserDefinedSchemaUsedByObject tracks when an object is created in a user defined schema.
	UserDefinedSchemaUsedByObject
)

var userDefinedSchemaTelemetryMap = map[UserDefinedSchemaTelemetryType]string{
	UserDefinedSchemaCreate:           "create_schema",
	UserDefinedSchemaAlter:            "alter_schema",
	UserDefinedSchemaDrop:             "drop_schema",
	UserDefinedSchemaReparentDatabase: "reparent_database",
	UserDefinedSchemaUsedByObject:     "schema_used_by_object",
}

func (u UserDefinedSchemaTelemetryType) String() string {
	return userDefinedSchemaTelemetryMap[u]
}

var userDefinedSchemaTelemetryCounters map[UserDefinedSchemaTelemetryType]telemetry.Counter

func init() {
	userDefinedSchemaTelemetryCounters = make(map[UserDefinedSchemaTelemetryType]telemetry.Counter)
	for ty, s := range userDefinedSchemaTelemetryMap {
		userDefinedSchemaTelemetryCounters[ty] = telemetry.GetCounterOnce(fmt.Sprintf("sql.uds.%s", s))
	}
}

// IncrementUserDefinedSchemaCounter is used to increment the telemetry counter
// for a particular usage of user defined schemas.
func IncrementUserDefinedSchemaCounter(userDefinedSchemaType UserDefinedSchemaTelemetryType) {
	telemetry.Inc(userDefinedSchemaTelemetryCounters[userDefinedSchemaType])
}
