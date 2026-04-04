// Copyright 2019 The Cockroach Authors.
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

package workloadimpl

import (
	"math"

	"golang.org/x/exp/rand"
)

// RandStringFast is a non-specialized random string generator with an even
// distribution of alphabet in the output.
func RandStringFast(rng rand.Source, buf []byte, alphabet string) {
	// We could pull these computations out to be done once per alphabet, but at
	// that point, you likely should consider PrecomputedRand.
	alphabetLen := uint64(len(alphabet))
	// floor(log(math.MaxUint64)/log(alphabetLen))
	lettersCharsPerRand := uint64(math.Log(float64(math.MaxUint64)) / math.Log(float64(alphabetLen)))

	var r, charsLeft uint64
	for i := 0; i < len(buf); i++ {
		if charsLeft == 0 {
			r = rng.Uint64()
			charsLeft = lettersCharsPerRand
		}
		buf[i] = alphabet[r%alphabetLen]
		r = r / alphabetLen
		charsLeft--
	}
}
