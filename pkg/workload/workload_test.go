// Copyright 2017 The Cockroach Authors.
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

package workload_test

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/workload"
)

func TestGet(t *testing.T) {
	defer leaktest.AfterTest(t)()
	if _, err := workload.Get(`bank`); err != nil {
		t.Errorf(`expected success got: %+v`, err)
	}
	if _, err := workload.Get(`nope`); !testutils.IsError(err, `unknown generator`) {
		t.Errorf(`expected "unknown generator" error got: %+v`, err)
	}
}

func TestApproxDatumSize(t *testing.T) {
	defer leaktest.AfterTest(t)()

	tests := []struct {
		datum interface{}
		size  int64
	}{
		{nil, 0},
		{int(-1000000), 3},
		{int(-1000), 2},
		{int(-1), 1},
		{int(0), 1},
		{int(1), 1},
		{int(1000), 2},
		{int(1000000), 3},
		{float64(0), 1},
		{float64(3), 8},
		{float64(3.1), 8},
		{float64(3.14), 8},
		{float64(3.141), 8},
		{float64(3.1415), 8},
		{float64(3.14159), 8},
		{float64(3.141592), 8},
		{"", 0},
		{"a", 1},
		{"aa", 2},
	}
	for _, test := range tests {
		if size := workload.ApproxDatumSize(test.datum); size != test.size {
			t.Errorf(`%T: %v got %d expected %d`, test.datum, test.datum, size, test.size)
		}
	}
}
