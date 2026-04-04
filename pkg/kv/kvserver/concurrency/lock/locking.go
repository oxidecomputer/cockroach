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

// Package lock provides type definitions for locking-related concepts used by
// concurrency control in the key-value layer.
package lock

import fmt "fmt"

// MaxDurability is the maximum value in the Durability enum.
const MaxDurability = Unreplicated

func init() {
	for v := range Durability_name {
		if d := Durability(v); d > MaxDurability {
			panic(fmt.Sprintf("Durability (%s) with value larger than MaxDurability", d))
		}
	}
}

// SafeValue implements redact.SafeValue.
func (Strength) SafeValue() {}

// SafeValue implements redact.SafeValue.
func (Durability) SafeValue() {}

// SafeValue implements redact.SafeValue.
func (WaitPolicy) SafeValue() {}
