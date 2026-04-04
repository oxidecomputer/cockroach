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

package faker

import "golang.org/x/exp/rand"

// This is a rough go port of https://github.com/joke2k/faker.

// Faker returns a generator for random data of various types: names, addresses,
// "lorem ipsum" placeholder text. Faker is safe for concurrent use.
type Faker struct {
	addressFaker
	loremFaker
	nameFaker
}

// NewFaker returns a new Faker instance. This causes a good amount of
// allocations, so it's best to reuse these when possible.
func NewFaker() Faker {
	f := Faker{
		loremFaker: newLoremFaker(),
		nameFaker:  newNameFaker(),
	}
	f.addressFaker = newAddressFaker(f.nameFaker)
	return f
}

type weightedEntry struct {
	weight float64
	entry  interface{}
}

type weightedEntries struct {
	entries     []weightedEntry
	totalWeight float64
}

func makeWeightedEntries(entriesAndWeights ...interface{}) *weightedEntries {
	we := make([]weightedEntry, 0, len(entriesAndWeights)/2)
	var totalWeight float64
	for idx := 0; idx < len(entriesAndWeights); idx += 2 {
		e, w := entriesAndWeights[idx], entriesAndWeights[idx+1].(float64)
		we = append(we, weightedEntry{weight: w, entry: e})
		totalWeight += w
	}
	return &weightedEntries{entries: we, totalWeight: totalWeight}
}

func (e *weightedEntries) Rand(rng *rand.Rand) interface{} {
	rngWeight := rng.Float64() * e.totalWeight
	var w float64
	for i := range e.entries {
		w += e.entries[i].weight
		if w > rngWeight {
			return e.entries[i].entry
		}
	}
	panic(`unreachable`)
}

func randInt(rng *rand.Rand, minInclusive, maxInclusive int) int {
	return rng.Intn(maxInclusive-minInclusive+1) + minInclusive
}
