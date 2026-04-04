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

package schemachange

import (
	"math/rand"

	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
)

// Deck is a random number generator that generates numbers in the range
// [0,len(weights)-1] where the probability of i is
// weights(i)/sum(weights). Unlike Weighted, the weights are specified as
// integers and used in a deck-of-cards style random number selection which
// ensures that each element is returned with a desired frequency within the
// size of the deck.
type deck struct {
	rng *rand.Rand
	mu  struct {
		syncutil.Mutex
		index int
		vals  []int
	}
}

// newDeck returns a new deck random number generator.
func newDeck(rng *rand.Rand, weights ...int) *deck {
	var sum int
	for i := range weights {
		sum += weights[i]
	}
	vals := make([]int, 0, sum)
	for i := range weights {
		for j := 0; j < weights[i]; j++ {
			vals = append(vals, i)
		}
	}
	d := &deck{
		rng: rng,
	}
	d.mu.index = len(vals)
	d.mu.vals = vals
	return d
}

// Int returns a random number in the range [0,len(weights)-1] where the
// probability of i is weights(i)/sum(weights).
func (d *deck) Int() int {
	d.mu.Lock()
	if d.mu.index == len(d.mu.vals) {
		d.rng.Shuffle(len(d.mu.vals), func(i, j int) {
			d.mu.vals[i], d.mu.vals[j] = d.mu.vals[j], d.mu.vals[i]
		})
		d.mu.index = 0
	}
	result := d.mu.vals[d.mu.index]
	d.mu.index++
	d.mu.Unlock()
	return result
}
