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
	"reflect"

	"github.com/cockroachdb/cockroach/pkg/util"
)

type slotIdx int

type fact struct {
	variable slotIdx
	attr     ordinal
	value    slotIdx
}

// slot represents an potentially unbound value referenced in a query.
type slot struct {
	typedValue

	// any holds the acceptable valuesMap which may occupy this slot as
	// indicated from an Any value.
	// TODO(ajwerner): Consider pulling this out into a set of constraints
	// on the slot. Then we could potentially also have the type act as a
	// constraint on the value. This mostly arises in the context of the
	// filter functions but conceivably could arise also when we have
	// equality filters and then variables used in attributes which require
	// types.
	any []typedValue
}

// typedValue is a value in its comparable form, which is to say, it is a
// pointer to a primitive type. The type is the type from which the value
// was derived and which should be presented to the user.
type typedValue struct {
	typ   reflect.Type
	value interface{}
}

func (tv typedValue) toInterface() interface{} {
	if tv.typ == reflectTypeType {
		return tv.value.(reflect.Type)
	}
	if tv.typ.Kind() == reflect.Ptr {
		if tv.typ.Elem().Kind() == reflect.Struct {
			return tv.value
		}
		return reflect.ValueOf(tv.value).Convert(tv.typ).Interface()
	}
	return reflect.ValueOf(tv.value).Convert(reflect.PtrTo(tv.typ)).Elem().Interface()
}

func (s *slot) eq(other slot) bool {
	// TODO(ajwerner): Deal with types. We may have two slots which both have
	// nil values but they differ in terms of types.
	switch {
	case s.value == nil && other.value == nil:
		return true
	case s.value == nil:
		return false
	case other.value == nil:
		return false
	default:
		_, eq := compare(s.value, other.value)
		return eq
	}
}

func (s *slot) empty() bool {
	return s.value == nil
}

func maybeSet(
	slots []slot, idx slotIdx, tv typedValue, set *util.FastIntSet,
) (foundContradiction bool) {
	s := &slots[idx]
	check := func() (shouldSet, foundContradiction bool) {
		if !s.empty() {
			if _, eq := compare(s.value, tv.value); !eq {
				return false, true
			}
			return false, false
		}

		if s.any != nil {
			var foundMatch bool
			for _, v := range s.any {
				if tv.typ != v.typ {
					continue
				}
				if _, foundMatch = compare(v.value, tv.value); foundMatch {
					break
				}
			}
			if !foundMatch {
				return false, true // contradiction
			}
		}
		return true, false
	}
	shouldSet, contradiction := check()
	if !shouldSet || contradiction {
		return contradiction
	}
	s.typedValue = tv
	if set != nil {
		set.Add(int(idx))
	}
	return false
}

// filter is a user-provided predicate over some set of variables.
type filter struct {
	input     []slotIdx
	predicate reflect.Value
}
