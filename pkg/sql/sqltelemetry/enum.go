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

// EnumTelemetryType represents a type of ENUM related operation to record
// telemetry for.
type EnumTelemetryType int

const (
	_ EnumTelemetryType = iota
	// EnumCreate represents a CREATE TYPE ... AS ENUM command.
	EnumCreate
	// EnumAlter represents an ALTER TYPE ... command.
	EnumAlter
	// EnumDrop represents a DROP TYPE command.
	EnumDrop
	// EnumInTable tracks when an enum type is used in a table.
	EnumInTable
)

var enumTelemetryMap = map[EnumTelemetryType]string{
	EnumCreate:  "create_enum",
	EnumAlter:   "alter_enum",
	EnumDrop:    "drop_enum",
	EnumInTable: "enum_used_in_table",
}

func (e EnumTelemetryType) String() string {
	return enumTelemetryMap[e]
}

var enumTelemetryCounters map[EnumTelemetryType]telemetry.Counter

func init() {
	enumTelemetryCounters = make(map[EnumTelemetryType]telemetry.Counter)
	for ty, s := range enumTelemetryMap {
		enumTelemetryCounters[ty] = telemetry.GetCounterOnce(fmt.Sprintf("sql.udts.%s", s))
	}
}

// IncrementEnumCounter is used to increment the telemetry counter for a particular
// usage of enums.
func IncrementEnumCounter(enumType EnumTelemetryType) {
	telemetry.Inc(enumTelemetryCounters[enumType])
}
