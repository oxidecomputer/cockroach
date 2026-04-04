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

package sessiondatapb

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/security"
)

// GetFloatPrec computes a precision suitable for a call to
// strconv.FormatFloat() or for use with '%.*g' in a printf-like
// function.
func (c DataConversionConfig) GetFloatPrec() int {
	// The user-settable parameter ExtraFloatDigits indicates the number
	// of digits to be used to format the float value. PostgreSQL
	// combines this with %g.
	// The formula is <type>_DIG + extra_float_digits,
	// where <type> is either FLT (float4) or DBL (float8).

	// Also the value "3" in PostgreSQL is special and meant to mean
	// "all the precision needed to reproduce the float exactly". The Go
	// formatter uses the special value -1 for this and activates a
	// separate path in the formatter. We compare >= 3 here
	// just in case the value is not gated properly in the implementation
	// of SET.
	if c.ExtraFloatDigits >= 3 {
		return -1
	}

	// CockroachDB only implements float8 at this time and Go does not
	// expose DBL_DIG, so we use the standard literal constant for
	// 64bit floats.
	const StdDoubleDigits = 15

	nDigits := StdDoubleDigits + c.ExtraFloatDigits
	if nDigits < 1 {
		// Ensure the value is clamped at 1: printf %g does not allow
		// values lower than 1. PostgreSQL does this too.
		nDigits = 1
	}
	return int(nDigits)
}

func (m VectorizeExecMode) String() string {
	switch m {
	case VectorizeOn, VectorizeUnset:
		return "on"
	case VectorizeExperimentalAlways:
		return "experimental_always"
	case VectorizeOff:
		return "off"
	default:
		return fmt.Sprintf("invalid (%d)", m)
	}
}

// VectorizeExecModeFromString converts a string into a VectorizeExecMode.
// False is returned if the conversion was unsuccessful.
func VectorizeExecModeFromString(val string) (VectorizeExecMode, bool) {
	var m VectorizeExecMode
	switch strings.ToUpper(val) {
	case "ON":
		m = VectorizeOn
	case "EXPERIMENTAL_ALWAYS":
		m = VectorizeExperimentalAlways
	case "OFF":
		m = VectorizeOff
	default:
		return 0, false
	}
	return m, true
}

// User retrieves the current user.
func (s *SessionData) User() security.SQLUsername {
	return s.UserProto.Decode()
}
