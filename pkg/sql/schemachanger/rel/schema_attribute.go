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

import (
	"unsafe"

	"github.com/cockroachdb/errors"
)

// GetAttribute returns the requested attribute from the passed entity. If
// the entity is of a type not defined for this schema, an error will be
// returned. If the entity is nil, an error will be returned. If the entity
// does not have this attribute defined, an error will be returned. If the
// entity does not have a populated value for this attribute, a nil value
// will be returned without an error.
func (sc *Schema) GetAttribute(attribute Attr, v interface{}) (interface{}, error) {
	ord, err := sc.getOrdinal(attribute)
	if err != nil {
		return nil, err
	}
	ti, value, err := getEntityValueInfo(sc, v)
	if err != nil {
		return nil, err
	}

	fi, ok := ti.attrFields[ord]
	if !ok {
		return nil, errors.Errorf(
			"no field defined on %v for %v", ti.typ, attribute,
		)
	}

	// There may be more than one field which stores this attribute. At most one
	// such field should be populated.
	for i := range fi {
		got := fi[i].value(unsafe.Pointer(value.Pointer()))
		if got != nil {
			return got, nil
		}
	}
	// No field was populated.
	return nil, nil
}
