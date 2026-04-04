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

package types

// The following variables are useful for testing.
var (
	// OneIntCol is a slice of one IntType.
	OneIntCol = []*T{Int}
	// TwoIntCols is a slice of two IntTypes.
	TwoIntCols = []*T{Int, Int}
	// ThreeIntCols is a slice of three IntTypes.
	ThreeIntCols = []*T{Int, Int, Int}
	// FourIntCols is a slice of four IntTypes.
	FourIntCols = []*T{Int, Int, Int, Int}
)

// MakeIntCols makes a slice of numCols IntTypes.
func MakeIntCols(numCols int) []*T {
	ret := make([]*T, numCols)
	for i := 0; i < numCols; i++ {
		ret[i] = Int
	}
	return ret
}
