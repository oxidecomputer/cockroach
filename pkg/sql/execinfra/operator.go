// Copyright 2018 The Cockroach Authors.
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

package execinfra

// OpNode is an interface to operator-like structures with children.
type OpNode interface {
	// ChildCount returns the number of children (inputs) of the operator.
	ChildCount(verbose bool) int

	// Child returns the nth child (input) of the operator.
	Child(nth int, verbose bool) OpNode
}

// OpChains describes a forest of OpNodes that represent a single physical plan.
// Each entry in the slice is a root of a separate OpNode tree.
type OpChains []OpNode
