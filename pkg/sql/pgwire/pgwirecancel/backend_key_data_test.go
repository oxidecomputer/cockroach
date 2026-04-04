// Copyright 2022 The Cockroach Authors.
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

package pgwirecancel

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/stretchr/testify/require"
)

func TestMakeBackendKeyData(t *testing.T) {
	defer leaktest.AfterTest(t)()

	rng := rand.New(rand.NewSource(timeutil.Now().Unix()))
	b1 := MakeBackendKeyData(rng, base.SQLInstanceID(1))
	b2 := MakeBackendKeyData(rng, base.SQLInstanceID(1))
	require.NotEqual(t, b1, b2)
}

func TestGetSQLInstanceID(t *testing.T) {
	defer leaktest.AfterTest(t)()

	rng := rand.New(rand.NewSource(timeutil.Now().Unix()))
	t.Run("small id", func(t *testing.T) {
		for i := 0; i < 1<<11; i++ {
			b := MakeBackendKeyData(rng, base.SQLInstanceID(i))
			require.Equal(t, base.SQLInstanceID(i), b.GetSQLInstanceID())
		}
	})

	for i := 1 << 11; i < math.MaxInt32; i = (i * 2) + 1 {
		t.Run(fmt.Sprintf("large id %d", i), func(t *testing.T) {
			b := MakeBackendKeyData(rng, base.SQLInstanceID(i))
			require.Equal(t, base.SQLInstanceID(i), b.GetSQLInstanceID())
		})
	}
}
