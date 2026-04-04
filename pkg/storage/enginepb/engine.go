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

package enginepb

import (
	"fmt"

	"github.com/cockroachdb/redact"
)

// Type implements the pflag.Value interface.
func (e *EngineType) Type() string { return "string" }

// String implements the pflag.Value interface.
func (e *EngineType) String() string { return redact.StringWithoutMarkers(e) }

// SafeFormat implements the refact.SafeFormatter interface.
func (e *EngineType) SafeFormat(p redact.SafePrinter, _ rune) {
	switch *e {
	case EngineTypeDefault:
		p.SafeString("default")
	case EngineTypePebble:
		p.SafeString("pebble")
	default:
		p.Printf("<unknown engine %d>", int32(*e))
	}
}

// Set implements the pflag.Value interface.
func (e *EngineType) Set(s string) error {
	switch s {
	case "default":
		*e = EngineTypeDefault
	case "pebble":
		*e = EngineTypePebble
	default:
		return fmt.Errorf("invalid storage engine: %s "+
			"(possible values: pebble)", s)
	}
	return nil
}
