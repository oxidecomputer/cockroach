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

package ring

import (
	"math/rand"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/stretchr/testify/require"
)

func TestRingBuffer(t *testing.T) {
	defer leaktest.AfterTest(t)()
	const operationCount = 100
	var buffer Buffer
	naiveBuffer := make([]interface{}, 0, operationCount)
	for i := 0; i < operationCount; i++ {
		switch rand.Intn(5) {
		case 0:
			buffer.AddFirst(i)
			naiveBuffer = append([]interface{}{i}, naiveBuffer...)
		case 1:
			buffer.AddLast(i)
			naiveBuffer = append(naiveBuffer, i)
		case 2:
			if len(naiveBuffer) > 0 {
				buffer.RemoveFirst()
				// NB: shift to preserve length.
				copy(naiveBuffer, naiveBuffer[1:])
				naiveBuffer = naiveBuffer[:len(naiveBuffer)-1]
			}
		case 3:
			if len(naiveBuffer) > 0 {
				buffer.RemoveLast()
				naiveBuffer = naiveBuffer[:len(naiveBuffer)-1]
			}
		case 4:
			// If there's extra capacity, resize to trim it.
			require.LessOrEqual(t, len(naiveBuffer), buffer.Cap())
			spareCap := buffer.Cap() - len(naiveBuffer)
			if spareCap > 0 {
				buffer.Resize(len(naiveBuffer) + rand.Intn(spareCap))
			}
		default:
			t.Fatal("unexpected")
		}

		require.Equal(t, naiveBuffer, buffer.all())
	}
}

func TestRingBufferCapacity(t *testing.T) {
	defer leaktest.AfterTest(t)()
	var b Buffer

	require.Panics(t, func() { b.Reserve(-1) })
	require.Equal(t, 0, b.Len())
	require.Equal(t, 0, b.Cap())

	b.Reserve(0)
	require.Equal(t, 0, b.Len())
	require.Equal(t, 0, b.Cap())

	b.AddFirst("a")
	require.Equal(t, 1, b.Len())
	require.Equal(t, 1, b.Cap())
	require.Panics(t, func() { b.Reserve(0) })
	require.Equal(t, 1, b.Len())
	require.Equal(t, 1, b.Cap())
	b.Reserve(1)
	require.Equal(t, 1, b.Len())
	require.Equal(t, 1, b.Cap())
	b.Reserve(2)
	require.Equal(t, 1, b.Len())
	require.Equal(t, 2, b.Cap())

	b.AddLast("z")
	require.Equal(t, 2, b.Len())
	require.Equal(t, 2, b.Cap())
	require.Panics(t, func() { b.Reserve(1) })
	require.Equal(t, 2, b.Len())
	require.Equal(t, 2, b.Cap())
	b.Reserve(2)
	require.Equal(t, 2, b.Len())
	require.Equal(t, 2, b.Cap())
	b.Reserve(9)
	require.Equal(t, 2, b.Len())
	require.Equal(t, 9, b.Cap())

	b.RemoveFirst()
	require.Equal(t, 1, b.Len())
	require.Equal(t, 9, b.Cap())
	b.Reserve(1)
	require.Equal(t, 1, b.Len())
	require.Equal(t, 9, b.Cap())
	b.RemoveLast()
	require.Equal(t, 0, b.Len())
	require.Equal(t, 9, b.Cap())
	b.Reserve(0)
	require.Equal(t, 0, b.Len())
	require.Equal(t, 9, b.Cap())

	b.Resize(3)
	require.Equal(t, 0, b.Len())
	require.Equal(t, 3, b.Cap())
}
