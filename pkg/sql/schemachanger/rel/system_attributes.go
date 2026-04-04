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

package rel

// systemAttribute is a type which represents attributes offered by the
// system for all entities stored in a database. In particular they capture
// the type and address of the variable. Unlike other attributes which only
// apply to entities, the systemAttributes Type and Self apply to all values.
//
// The system attribute may be extended to cover other structural attributes.
// TODO(ajwerner): Add support for slices, arrays, and maps and then provide
// system attributes to access slice/array indexes and map keys and valuesMap.
type systemAttribute int8

//go:generate stringer -type systemAttribute

const (
	_ systemAttribute = 64 - iota

	// Type is an attribute which stores a value's type.
	Type

	// Self is an attribute which stores the variable itself.
	Self

	maxUserAttribute ordinal = 64 - iota
)

func isSystemAttribute(a Attr) bool {
	_, isSystemAttr := a.(systemAttribute)
	return isSystemAttr
}

var _ Attr = systemAttribute(0)
