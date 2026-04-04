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

package optional

import "time"

// MakeTimeValue returns an Duration with a set value.
func MakeTimeValue(value time.Duration) Duration {
	return Duration{ValuePlusOne: value + 1}
}

// HasValue returns true if a value was set.
func (d Duration) HasValue() bool {
	return d.ValuePlusOne != 0
}

// Value returns the current value, or 0 if HasValue() is false.
func (d Duration) Value() time.Duration {
	if d.ValuePlusOne == 0 {
		return 0
	}
	return d.ValuePlusOne - 1
}

func (d Duration) String() string {
	if !d.HasValue() {
		return "<unset>"
	}
	return d.Value().String()
}

// Clear the value.
func (d *Duration) Clear() {
	d.ValuePlusOne = 0
}

// Set the value.
func (d *Duration) Set(value time.Duration) {
	*d = MakeTimeValue(value)
}

// Add modifies the value by adding a delta.
func (d *Duration) Add(delta time.Duration) {
	*d = MakeTimeValue(d.Value() + delta)
}

// MaybeAdd adds the given value, if it is set. Does nothing if other is not
// set.
func (d *Duration) MaybeAdd(other Duration) {
	if other.HasValue() {
		*d = MakeTimeValue(d.Value() + other.Value())
	}
}
