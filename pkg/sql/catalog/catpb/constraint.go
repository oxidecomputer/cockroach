// Copyright 2022 The Cockroach Authors.
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

package catpb

import (
	"strconv"

	"github.com/cockroachdb/redact"
)

// String implements the fmt.Stringer interface.
func (x ForeignKeyAction) String() string {
	switch x {
	case ForeignKeyAction_RESTRICT:
		return "RESTRICT"
	case ForeignKeyAction_SET_DEFAULT:
		return "SET DEFAULT"
	case ForeignKeyAction_SET_NULL:
		return "SET NULL"
	case ForeignKeyAction_CASCADE:
		return "CASCADE"
	default:
		return strconv.Itoa(int(x))
	}
}

var _ redact.SafeValue = ForeignKeyAction(0)

// SafeValue implements redact.SafeValue.
func (x ForeignKeyAction) SafeValue() {}
