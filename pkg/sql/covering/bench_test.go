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

package covering

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkOverlapCoveringMerge(b *testing.B) {
	var benchmark []struct {
		name   string
		inputs []Covering
	}

	for _, numLayers := range []int{
		1,      // single backup
		24,     // hourly backups
		24 * 7, // hourly backups for a week.
	} {
		// number of elements per each backup instance
		for _, elementsPerLayer := range []int{100, 1000, 10000} {
			var inputs []Covering

			for i := 0; i < numLayers; i++ {
				var payload int
				var c Covering
				step := 1 + rand.Intn(10)

				for j := 0; j < elementsPerLayer; j += step {
					start := make([]byte, 4)
					binary.LittleEndian.PutUint32(start, uint32(j))

					end := make([]byte, 4)
					binary.LittleEndian.PutUint32(end, uint32(j+step))

					c = append(c, Range{
						Start:   start,
						End:     end,
						Payload: payload,
					})
					payload++
				}
				inputs = append(inputs, c)
			}

			benchmark = append(benchmark, struct {
				name   string
				inputs []Covering
			}{name: fmt.Sprintf("layers=%d,elems=%d", numLayers, elementsPerLayer), inputs: inputs})
		}
	}

	b.ResetTimer()
	for _, bench := range benchmark {
		inputs := bench.inputs
		b.Run(bench.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				require.NotEmpty(b, OverlapCoveringMerge(inputs))
			}
		})
	}
}
