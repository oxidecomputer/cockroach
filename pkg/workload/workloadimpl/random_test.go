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

package workloadimpl_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/cockroach/pkg/workload/workloadimpl"
	"golang.org/x/exp/rand"
)

func BenchmarkRandStringFast(b *testing.B) {
	const strLen = 26
	rng := rand.NewSource(uint64(timeutil.Now().UnixNano()))
	buf := make([]byte, strLen)

	for i := 0; i < b.N; i++ {
		workloadimpl.RandStringFast(rng, buf, `0123456789abcdefghijklmnopqrstuvwxyz`)
	}
	b.SetBytes(strLen)
}
