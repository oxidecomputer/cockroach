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

package opt

import (
	"reflect"
	"runtime"
	"testing"
)

// TestAggregateProperties verifies that the various helper functions for
// various properties of aggregations handle all aggregation operators.
func TestAggregateProperties(t *testing.T) {
	check := func(fn func()) bool {
		ok := true
		func() {
			defer func() {
				if x := recover(); x != nil {
					ok = false
				}
			}()
			fn()
		}()
		return ok
	}

	for _, op := range AggregateOperators {
		funcs := []func(Operator) bool{
			AggregateIgnoresDuplicates,
			AggregateIgnoresNulls,
			AggregateIsNeverNull,
			AggregateIsNeverNullOnNonNullInput,
			AggregateIsNullOnEmpty,
		}

		for _, fn := range funcs {
			if !check(func() { fn(op) }) {
				fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				t.Errorf("%s not handled by %s", op, fnName)
			}
		}

		for _, op2 := range AggregateOperators {
			if !check(func() { AggregatesCanMerge(op, op2) }) {
				t.Errorf("%s,%s not handled by AggregatesCanMerge", op, op2)
				break
			}
		}
	}
}
