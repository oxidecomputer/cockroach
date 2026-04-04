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

package kvnemesis

import (
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/assert"
)

func TestOperationsFormat(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	tests := []struct {
		step     Step
		expected string
	}{
		{step: step(get(`a`)), expected: `db0.Get(ctx, "a")`},
		{step: step(del(`a`)), expected: `db0.Del(ctx, "a")`},
		{step: step(batch(get(`b`), reverseScanForUpdate(`c`, `e`), get(`f`))), expected: `
			{
			  b := &Batch{}
			  b.Get(ctx, "b")
			  b.ReverseScanForUpdate(ctx, "c", "e")
			  b.Get(ctx, "f")
			  db0.Run(ctx, b)
			}
		`},
		{
			step: step(
				closureTxn(ClosureTxnType_Commit,
					batch(get(`g`), get(`h`), del(`i`)),
					delRange(`j`, `k`),
					put(`k`, `l`),
				)),
			expected: `
			db0.Txn(ctx, func(ctx context.Context, txn *kv.Txn) error {
			  {
			    b := &Batch{}
			    b.Get(ctx, "g")
			    b.Get(ctx, "h")
			    b.Del(ctx, "i")
			    txn.Run(ctx, b)
			  }
			  txn.DelRange(ctx, "j", "k", true)
			  txn.Put(ctx, "k", l)
			  return nil
			})
			`,
		},
	}

	for _, test := range tests {
		expected := strings.TrimSpace(test.expected)
		var actual strings.Builder
		test.step.format(&actual, formatCtx{indent: "\t\t\t"})
		assert.Equal(t, expected, strings.TrimSpace(actual.String()))
	}
}
